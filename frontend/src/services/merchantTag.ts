import AbstractService from './abstractService'
import MerchantTagModel from '@/models/merchantTag'
import type {IMerchantTag} from '@/modelTypes/IMerchantTag'

export default class MerchantTagService extends AbstractService<IMerchantTag> {
	constructor() {
		super({
			create: '/merchant-tags',
			get: '/merchant-tags/{id}',
			getAll: '/merchant-tags',
			update: '/merchant-tags/{id}',
			delete: '/merchant-tags/{id}',
		})
	}

	modelFactory(data: Partial<IMerchantTag>) {
		return new MerchantTagModel(data)
	}

	modelCreateFactory(data: Partial<IMerchantTag>) {
		return new MerchantTagModel(data)
	}

	modelUpdateFactory(data: Partial<IMerchantTag>) {
		return new MerchantTagModel(data)
	}

	modelGetFactory(data: Partial<IMerchantTag>) {
		return new MerchantTagModel(data)
	}

	modelGetAllFactory(data: Partial<IMerchantTag>) {
		return new MerchantTagModel(data)
	}

	async getDistinctClasses() {
		const cancel = this.setLoading()
		try {
			const response = await this.http.get('/merchant-tags/classes')
			return response.data as string[]
		} finally {
			cancel()
		}
	}
}
