import type {IAbstract} from './IAbstract'
import type {IUser} from './IUser'
import type {ISubscription} from './ISubscription'
import type {IMerchantTag} from './IMerchantTag'
import type {IGeoPoint} from './IGeoPoint'

export interface IMerchant extends IAbstract {
	id: number
	title: string
	description: string
	legalName: string
	phone: string
	address: string
	city: string
	area: string
	
	// Geographical information
	lng: number | null
	lat: number | null
	geocodeLevel: string
	geocodeScore: number
	geocodeDescription: string
	geocodeAttempts: number
	
	// Relationships
	owner: IUser
	createdBy: IUser
	tags: IMerchantTag[]
	geoPoints: IGeoPoint[]
	
	// UI state
	isFavorite: boolean
	subscription: ISubscription | null
	
	// Timestamps
	created: Date
	updated: Date
}
