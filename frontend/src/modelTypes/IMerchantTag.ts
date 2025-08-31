import type {IAbstract} from './IAbstract'
import type {IUser} from './IUser'

export interface IMerchantTag extends IAbstract {
	id: number
	tagName: string
	alias: string
	class: string
	remarks: string
	hexColor: string
	
	// Relationships
	owner: IUser
	
	// Timestamps
	created: Date
	updated: Date
}
