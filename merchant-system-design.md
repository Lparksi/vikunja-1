# 商户系统迁移到Vikunja架构设计文档

## 1. 系统架构对比分析

### 参考系统特点
- **数据库**: GORM + PostgreSQL + Redis缓存
- **特性**: 地理编码、位置信息处理、Excel导入、商户标签分类管理
- **架构**: Repository模式 + 事务管理

### Vikunja系统特点  
- **数据库**: XORM + 多数据库支持(MySQL/PostgreSQL/SQLite)
- **权限系统**: 三层权限(Read/Write/Admin)
- **架构**: 事件驱动 + 标准化CRUD + web.CRUDable接口

## 2. 数据模型设计

### 2.1 商户模型 (Merchant)

```go
// Merchant represents a merchant entity
type Merchant struct {
	// 基础字段
	ID          int64     `xorm:"bigint autoincr not null unique pk" json:"id"`
	LegalName   string    `xorm:"varchar(100) not null" json:"legal_name" valid:"required,runelength(1|100)"`
	Phone       string    `xorm:"varchar(20)" json:"phone" valid:"runelength(0|20)"`
	Address     string    `xorm:"longtext" json:"address"`
	City        string    `xorm:"varchar(50)" json:"city" valid:"runelength(0|50)"`
	Area        string    `xorm:"varchar(50)" json:"area" valid:"runelength(0|50)"`
	
	// 地理编码信息
	Lng                *float64 `xorm:"decimal(11,8)" json:"lng"`
	Lat                *float64 `xorm:"decimal(10,8)" json:"lat"`
	GeocodeLevel       string   `xorm:"varchar(50)" json:"geocode_level"`
	GeocodeScore       int      `xorm:"int default 0" json:"geocode_score"`
	GeocodeDescription string   `xorm:"varchar(100) default '等待解析'" json:"geocode_description"`
	GeocodeAttempts    int      `xorm:"int default 0" json:"geocode_attempts"`
	
	// 关联信息
	OwnerID    int64        `xorm:"bigint INDEX not null" json:"-"`
	Owner      *user.User   `xorm:"-" json:"owner" valid:"-"`
	Tags       []*MerchantTag `xorm:"-" json:"tags"`
	
	// 地理位置点
	ActiveGeoPointID *int64     `xorm:"bigint" json:"active_geo_point_id"`
	ActiveGeoPoint   *GeoPoint  `xorm:"-" json:"active_geo_point"`
	GeoPoints        []*GeoPoint `xorm:"-" json:"geo_points"`
	
	// 收藏和权限
	IsFavorite    bool       `xorm:"-" json:"is_favorite"`
	Subscription  *Subscription `xorm:"-" json:"subscription,omitempty"`
	MaxPermission Permission `xorm:"-" json:"max_permission"`
	
	// 时间戳
	Created time.Time `xorm:"created not null" json:"created"`
	Updated time.Time `xorm:"updated not null" json:"updated"`
	
	// 实现接口
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}
```

### 2.2 商户标签模型 (MerchantTag)

```go
// MerchantTag represents a merchant tag
type MerchantTag struct {
	ID      int64  `xorm:"bigint autoincr not null unique pk" json:"id"`
	TagName string `xorm:"varchar(50) not null" json:"tag_name" valid:"required,runelength(1|50)"`
	Alias   string `xorm:"varchar(10)" json:"alias" valid:"runelength(0|10)"`
	Class   string `xorm:"varchar(50)" json:"class" valid:"runelength(0|50)"`
	Remarks string `xorm:"longtext" json:"remarks"`
	Color   string `xorm:"varchar(6)" json:"color" valid:"runelength(0|7)"`
	
	// 权限和所有者
	OwnerID       int64      `xorm:"bigint INDEX not null" json:"-"`
	Owner         *user.User `xorm:"-" json:"owner" valid:"-"`
	MaxPermission Permission `xorm:"-" json:"max_permission"`
	
	// 时间戳
	Created time.Time `xorm:"created not null" json:"created"`
	Updated time.Time `xorm:"updated not null" json:"updated"`
	
	// 实现接口
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}
```

### 2.3 地理位置点模型 (GeoPoint)

```go
// GeoPoint represents a geographical point
type GeoPoint struct {
	ID         int64   `xorm:"bigint autoincr not null unique pk" json:"id"`
	MerchantID int64   `xorm:"bigint INDEX not null" json:"merchant_id"`
	From       string  `xorm:"varchar(50)" json:"from"`
	Longitude  float64 `xorm:"decimal(11,8)" json:"longitude"`
	Latitude   float64 `xorm:"decimal(10,8)" json:"latitude"`
	Address    string  `xorm:"longtext" json:"address"`
	Accuracy   int     `xorm:"int default 0" json:"accuracy"`
	
	// 几何数据 (适配不同数据库)
	GeomData string `xorm:"longtext" json:"geom_data"` // WKT格式存储
	
	// 时间戳
	Created time.Time `xorm:"created not null" json:"created"`
	Updated time.Time `xorm:"updated not null" json:"updated"`
	
	// 实现接口
	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}
```

### 2.4 商户-标签关联表

```go
// MerchantTagRelation represents the relationship between merchants and tags
type MerchantTagRelation struct {
	ID           int64 `xorm:"bigint autoincr not null unique pk"`
	MerchantID   int64 `xorm:"bigint INDEX not null"`
	MerchantTagID int64 `xorm:"bigint INDEX not null"`
	Created      time.Time `xorm:"created not null"`
}
```

## 3. 权限系统集成

### 3.1 权限级别定义
```go
const (
	PermissionRead  Permission = 1 // 查看商户信息
	PermissionWrite Permission = 2 // 编辑商户信息、添加标签
	PermissionAdmin Permission = 3 // 删除商户、管理权限、Excel导入
)
```

### 3.2 权限检查实现
```go
// CanRead checks if the user can read the merchant
func (m *Merchant) CanRead(s *xorm.Session, a web.Auth) (bool, error) {
	// 所有者可以读取
	if m.OwnerID == a.GetID() {
		return true, nil
	}
	
	// 检查共享权限
	return checkMerchantPermissions(s, m.ID, a, PermissionRead)
}

// CanWrite checks if the user can write to the merchant
func (m *Merchant) CanWrite(s *xorm.Session, a web.Auth) (bool, error) {
	// 所有者可以写入
	if m.OwnerID == a.GetID() {
		return true, nil
	}
	
	// 检查共享权限
	return checkMerchantPermissions(s, m.ID, a, PermissionWrite)
}

// CanDelete checks if the user can delete the merchant
func (m *Merchant) CanDelete(s *xorm.Session, a web.Auth) (bool, error) {
	// 只有所有者或管理员权限可以删除
	if m.OwnerID == a.GetID() {
		return true, nil
	}
	
	return checkMerchantPermissions(s, m.ID, a, PermissionAdmin)
}
```

## 4. API路由设计

### 4.1 商户管理API
```
GET    /api/v1/merchants                    # 获取商户列表
POST   /api/v1/merchants                    # 创建商户
GET    /api/v1/merchants/{id}               # 获取商户详情
PUT    /api/v1/merchants/{id}               # 更新商户
DELETE /api/v1/merchants/{id}               # 删除商户

POST   /api/v1/merchants/import/excel       # Excel批量导入
GET    /api/v1/merchants/{id}/geocode       # 地理编码
POST   /api/v1/merchants/{id}/geocode       # 触发地理编码

GET    /api/v1/merchants/{id}/geopoints     # 获取地理位置点
POST   /api/v1/merchants/{id}/geopoints     # 创建地理位置点
PUT    /api/v1/geopoints/{id}               # 更新地理位置点
DELETE /api/v1/geopoints/{id}               # 删除地理位置点
```

### 4.2 标签管理API
```
GET    /api/v1/merchant-tags                # 获取标签列表
POST   /api/v1/merchant-tags                # 创建标签
GET    /api/v1/merchant-tags/{id}           # 获取标签详情
PUT    /api/v1/merchant-tags/{id}           # 更新标签
DELETE /api/v1/merchant-tags/{id}           # 删除标签

GET    /api/v1/merchant-tags/classes        # 获取标签分类
```

## 5. 前端页面架构

### 5.1 路由设计
```javascript
const routes = [
  {
    path: '/merchants',
    name: 'merchants.index',
    component: () => import('@/views/merchants/Index.vue'),
    meta: { title: 'merchant.title' }
  },
  {
    path: '/merchants/new',
    name: 'merchants.create',
    component: () => import('@/views/merchants/NewMerchant.vue'),
    meta: { title: 'merchant.create.title' }
  },
  {
    path: '/merchants/:id/edit',
    name: 'merchants.edit',
    component: () => import('@/views/merchants/EditMerchant.vue'),
    meta: { title: 'merchant.edit.title' }
  },
  {
    path: '/merchants/import',
    name: 'merchants.import',
    component: () => import('@/views/merchants/ImportMerchants.vue'),
    meta: { title: 'merchant.import.title' }
  }
]
```

### 5.2 组件架构
```
merchants/
├── Index.vue                 # 商户列表主页
├── NewMerchant.vue          # 新建商户
├── EditMerchant.vue         # 编辑商户
├── ImportMerchants.vue      # Excel导入
├── components/
│   ├── MerchantList.vue     # 商户列表组件
│   ├── MerchantCard.vue     # 商户卡片
│   ├── MerchantMap.vue      # 地图显示
│   ├── MerchantForm.vue     # 商户表单
│   ├── TagSelector.vue      # 标签选择器
│   ├── GeocodingPanel.vue   # 地理编码面板
│   └── ExcelUpload.vue      # Excel上传组件
```

## 6. 数据库迁移

### 6.1 创建商户表
```sql
CREATE TABLE merchants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    legal_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address LONGTEXT,
    city VARCHAR(50),
    area VARCHAR(50),
    lng DECIMAL(11,8),
    lat DECIMAL(10,8),
    geocode_level VARCHAR(50),
    geocode_score INT DEFAULT 0,
    geocode_description VARCHAR(100) DEFAULT '等待解析',
    geocode_attempts INT DEFAULT 0,
    owner_id BIGINT NOT NULL,
    active_geo_point_id BIGINT,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    INDEX idx_owner_id (owner_id),
    INDEX idx_geo_location (lng, lat)
);
```

### 6.2 创建商户标签表
```sql
CREATE TABLE merchant_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tag_name VARCHAR(50) NOT NULL,
    alias VARCHAR(10),
    class VARCHAR(50),
    remarks LONGTEXT,
    color VARCHAR(6),
    owner_id BIGINT NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    INDEX idx_owner_id (owner_id),
    INDEX idx_class (class)
);
```

### 6.3 创建地理位置点表
```sql
CREATE TABLE geo_points (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    from_source VARCHAR(50),
    longitude DECIMAL(11,8),
    latitude DECIMAL(10,8),
    address LONGTEXT,
    accuracy INT DEFAULT 0,
    geom_data LONGTEXT,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    INDEX idx_merchant_id (merchant_id),
    INDEX idx_geo_location (longitude, latitude)
);
```

## 7. 服务层设计

### 7.1 地理编码服务
```go
// GeocodeService handles geocoding operations
type GeocodeService struct {
	config *config.GeocodeConfig
}

func (gs *GeocodeService) GeocodeAddress(address string) (*GeocodeResult, error) {
	// 支持多种地理编码提供商
	// - Google Maps API
	// - 百度地图API
	// - 高德地图API
	// - OpenStreetMap Nominatim
}

func (gs *GeocodeService) BatchGeocode(addresses []string) ([]*GeocodeResult, error) {
	// 批量地理编码处理
}
```

### 7.2 Excel导入服务
```go
// ExcelImportService handles Excel file imports
type ExcelImportService struct {
	merchantService *MerchantService
}

func (eis *ExcelImportService) ImportFromExcel(file io.Reader, userID int64) (*ImportResult, error) {
	// Excel文件解析
	// 数据验证
	// 批量导入
	// 错误处理和报告
}
```

## 8. 国际化支持

### 8.1 中文语言包 (zh-CN.json)
```json
{
  "merchant": {
    "title": "商户管理",
    "create": {
      "title": "新建商户",
      "success": "商户创建成功"
    },
    "edit": {
      "title": "编辑商户",
      "success": "商户更新成功"
    },
    "fields": {
      "legalName": "法人姓名",
      "phone": "联系电话",
      "address": "地址",
      "city": "城市",
      "area": "区域",
      "tags": "标签"
    },
    "import": {
      "title": "批量导入",
      "uploadFile": "上传Excel文件",
      "processing": "正在处理...",
      "success": "导入完成",
      "errors": "导入错误"
    },
    "geocoding": {
      "title": "地理编码",
      "pending": "等待解析",
      "success": "解析成功",
      "failed": "解析失败",
      "retry": "重新解析"
    }
  },
  "merchantTag": {
    "title": "商户标签",
    "create": "创建标签",
    "fields": {
      "tagName": "标签名称",
      "alias": "别名",
      "class": "分类",
      "color": "颜色"
    }
  }
}
```

## 9. 事件系统集成

### 9.1 商户事件定义
```go
// MerchantCreatedEvent is fired when a merchant is created
type MerchantCreatedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// MerchantUpdatedEvent is fired when a merchant is updated  
type MerchantUpdatedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// MerchantDeletedEvent is fired when a merchant is deleted
type MerchantDeletedEvent struct {
	Merchant *Merchant
	Doer     web.Auth
}

// MerchantGeocodedEvent is fired when a merchant is geocoded
type MerchantGeocodedEvent struct {
	Merchant *Merchant
	Result   *GeocodeResult
}
```

## 10. 测试策略

### 10.1 单元测试
- 模型验证测试
- CRUD操作测试  
- 权限检查测试
- 地理编码服务测试

### 10.2 集成测试
- API端点测试
- 数据库迁移测试
- Excel导入功能测试

### 10.3 前端测试
- 组件单元测试
- 页面集成测试
- E2E测试

## 11. 性能优化

### 11.1 数据库优化
- 适当的索引设计
- 地理查询优化
- 分页查询优化

### 11.2 缓存策略
- 商户列表缓存
- 地理编码结果缓存
- 标签数据缓存

### 11.3 前端优化
- 懒加载
- 虚拟滚动
- 地图数据分块加载

## 12. 部署注意事项

### 12.1 环境配置
```yaml
geocoding:
  provider: "google" # google, baidu, amap, nominatim
  api_key: "your-api-key"
  rate_limit: 100
  timeout: 30s

merchants:
  excel_import:
    max_file_size: "10MB"
    batch_size: 1000
    allowed_extensions: [".xlsx", ".xls"]
```

### 12.2 数据库支持
- MySQL: 完整功能支持
- PostgreSQL: 完整功能支持 + PostGIS扩展
- SQLite: 基础功能支持（地理数据以文本存储）

这个设计确保了商户系统能够完美融入Vikunja的架构中，保持一致的代码风格、权限模型和用户体验。
