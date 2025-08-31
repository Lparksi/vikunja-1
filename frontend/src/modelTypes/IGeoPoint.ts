import type {IAbstract} from './IAbstract'

export interface IGeoMetadata {
	provider?: string
	confidence?: number
	placeId?: string
	placeType?: string
	category?: string
}

export interface IGeoPoint extends IAbstract {
	id: number
	merchantId: number
	from: string
	longitude: number
	latitude: number
	address: string
	accuracy: number
	metadata: IGeoMetadata | null
	
	// Timestamps
	created: Date
	updated: Date
}
