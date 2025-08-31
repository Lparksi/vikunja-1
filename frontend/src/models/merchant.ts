import AbstractModel from './abstractModel'
import UserModel from '@/models/user'
import SubscriptionModel from '@/models/subscription'
import MerchantTagModel from '@/models/merchantTag'
import GeoPointModel from '@/models/geoPoint'

import type {IMerchant} from '@/modelTypes/IMerchant'
import type {IUser} from '@/modelTypes/IUser'
import type {ISubscription} from '@/modelTypes/ISubscription'
import type {IMerchantTag} from '@/modelTypes/IMerchantTag'
import type {IGeoPoint} from '@/modelTypes/IGeoPoint'

export default class MerchantModel extends AbstractModel<IMerchant> implements IMerchant {
	id = 0
	title = ''
	description = ''
	legalName = ''
	phone = ''
	address = ''
	city = ''
	area = ''
	
	// Geographical information
	lng: number | null = null
	lat: number | null = null
	geocodeLevel = ''
	geocodeScore = 0
	geocodeDescription = '等待解析'
	geocodeAttempts = 0
	
	// Relationships
	owner: IUser = new UserModel()
	createdBy: IUser = new UserModel()
	tags: IMerchantTag[] = []
	geoPoints: IGeoPoint[] = []
	
	// UI state
	isFavorite = false
	subscription: ISubscription | null = null
	
	// Timestamps
	created: Date = new Date()
	updated: Date = new Date()

	constructor(data: Partial<IMerchant> = {}) {
		super()
		this.assignData(data)

		this.owner = new UserModel(this.owner)
		this.createdBy = new UserModel(this.createdBy)

		// Make all tags to tag models
		this.tags = this.tags.map(t => {
			return new MerchantTagModel(t)
		})

		// Make all geo points to geo point models
		this.geoPoints = this.geoPoints.map(gp => {
			return new GeoPointModel(gp)
		})

		if (typeof this.subscription !== 'undefined' && this.subscription !== null) {
			this.subscription = new SubscriptionModel(this.subscription)
		}
		
		this.created = new Date(this.created)
		this.updated = new Date(this.updated)
	}

	// Computed properties
	get hasLocation(): boolean {
		return this.lng !== null && this.lat !== null
	}

	get isGeocoded(): boolean {
		return this.geocodeDescription !== '等待解析' && this.hasLocation
	}

	get fullAddress(): string {
		const parts = [this.address, this.area, this.city].filter(Boolean)
		return parts.join(', ')
	}
}
