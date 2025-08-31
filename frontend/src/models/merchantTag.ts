import AbstractModel from './abstractModel'
import UserModel from '@/models/user'

import type {IMerchantTag} from '@/modelTypes/IMerchantTag'
import type {IUser} from '@/modelTypes/IUser'

export default class MerchantTagModel extends AbstractModel<IMerchantTag> implements IMerchantTag {
	id = 0
	tagName = ''
	alias = ''
	class = ''
	remarks = ''
	hexColor = ''
	
	// Relationships
	owner: IUser = new UserModel()
	
	// Timestamps
	created: Date = new Date()
	updated: Date = new Date()

	constructor(data: Partial<IMerchantTag> = {}) {
		super()
		this.assignData(data)

		this.owner = new UserModel(this.owner)

		if (this.hexColor !== '' && this.hexColor.substring(0, 1) !== '#') {
			this.hexColor = '#' + this.hexColor
		}
		
		this.created = new Date(this.created)
		this.updated = new Date(this.updated)
	}
}
