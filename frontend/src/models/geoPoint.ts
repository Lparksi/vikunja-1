import AbstractModel from './abstractModel'

import type {IGeoPoint, IGeoMetadata} from '@/modelTypes/IGeoPoint'

export default class GeoPointModel extends AbstractModel<IGeoPoint> implements IGeoPoint {
	id = 0
	merchantId = 0
	from = ''
	longitude = 0
	latitude = 0
	address = ''
	accuracy = 0
	metadata: IGeoMetadata | null = null
	
	// Timestamps
	created: Date = new Date()
	updated: Date = new Date()

	constructor(data: Partial<IGeoPoint> = {}) {
		super()
		this.assignData(data)
		
		this.created = new Date(this.created)
		this.updated = new Date(this.updated)
	}

	// Computed properties
	get coordinates(): [number, number] {
		return [this.longitude, this.latitude]
	}

	get isValid(): boolean {
		return this.longitude >= -180 && this.longitude <= 180 &&
			   this.latitude >= -90 && this.latitude <= 90
	}
}
