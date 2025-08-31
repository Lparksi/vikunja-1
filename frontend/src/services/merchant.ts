import AbstractService from './abstractService'
import MerchantModel from '@/models/merchant'
import type {IMerchant} from '@/modelTypes/IMerchant'

export default class MerchantService extends AbstractService<IMerchant> {
	constructor() {
		super({
			create: '/merchants',
			get: '/merchants/{id}',
			getAll: '/merchants',
			update: '/merchants/{id}',
			delete: '/merchants/{id}',
		})
	}

	modelFactory(data: Partial<IMerchant>) {
		return new MerchantModel(data)
	}

	modelCreateFactory(data: Partial<IMerchant>) {
		return new MerchantModel(data)
	}

	modelUpdateFactory(data: Partial<IMerchant>) {
		return new MerchantModel(data)
	}

	modelGetFactory(data: Partial<IMerchant>) {
		return new MerchantModel(data)
	}

	modelGetAllFactory(data: Partial<IMerchant>) {
		return new MerchantModel(data)
	}

	async geocodeMerchant(merchant: IMerchant) {
		const cancel = this.setLoading()
		try {
			const response = await this.http.post(`/merchants/${merchant.id}/geocode`)
			return this.modelUpdateFactory(response.data)
		} finally {
			cancel()
		}
	}

	async importFromExcel(file: File) {
		const cancel = this.setLoading()
		try {
			const formData = new FormData()
			formData.append('file', file)
			
			const response = await this.http.post('/merchants/import/excel', formData, {
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				onUploadProgress: ({progress}) => {
					this.uploadProgress = progress ? Math.round((progress * 100)) : 0
				},
			})
			
			return response.data
		} finally {
			this.uploadProgress = 0
			cancel()
		}
	}

	async importFromCSV(file: File, autoGeocode: boolean = true) {
		const cancel = this.setLoading()
		try {
			const formData = new FormData()
			formData.append('file', file)
			formData.append('autoGeocode', autoGeocode.toString())

			const response = await this.http.post('/merchants/import/csv', formData, {
				headers: {
					'Content-Type': 'multipart/form-data',
				},
				onUploadProgress: ({progress}) => {
					this.uploadProgress = progress ? Math.round((progress * 100)) : 0
				},
			})
			
			return response.data
		} finally {
			this.uploadProgress = 0
			cancel()
		}
	}

	async geocode(address: string) {
		const cancel = this.setLoading()
		try {
			const response = await this.http.post('/merchants/geocode', {address})
			return response.data
		} finally {
			cancel()
		}
	}
}
