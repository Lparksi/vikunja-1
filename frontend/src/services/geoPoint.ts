import AbstractService from './abstractService'
import GeoPointModel from '@/models/geoPoint'
import type {IGeoPoint} from '@/modelTypes/IGeoPoint'

export default class GeoPointService extends AbstractService<IGeoPoint> {
	constructor() {
		super({
			create: '/merchants/{merchantId}/geopoints',
			get: '/geopoints/{id}',
			getAll: '/merchants/{merchantId}/geopoints',
			update: '/geopoints/{id}',
			delete: '/geopoints/{id}',
		})
	}

	modelFactory(data: Partial<IGeoPoint>) {
		return new GeoPointModel(data)
	}

	modelCreateFactory(data: Partial<IGeoPoint>) {
		return new GeoPointModel(data)
	}

	modelUpdateFactory(data: Partial<IGeoPoint>) {
		return new GeoPointModel(data)
	}

	modelGetFactory(data: Partial<IGeoPoint>) {
		return new GeoPointModel(data)
	}

	modelGetAllFactory(data: Partial<IGeoPoint>) {
		return new GeoPointModel(data)
	}
}
