<template>
	<div class="loader-container is-max-width-desktop" :class="{ 'is-loading': loading }">
		<div class="card">
			<header class="card-header">
				<p class="card-header-title">
					{{ $t('merchant.edit.title') }}
				</p>
			</header>
			<div class="card-content">
				<div class="content">
					<form @submit.prevent="updateMerchant()">
						<div class="field">
							<label class="label" for="merchantTitle">
								{{ $t('merchant.attributes.title') }}
								<span class="is-required">*</span>
							</label>
							<div class="control">
								<input
									id="merchantTitle"
									v-model="merchant.title"
									v-focus
									:class="{ 'is-danger': titleError }"
									class="input"
									:placeholder="$t('merchant.edit.titlePlaceholder')"
									type="text"
									required
									maxlength="100"
								>
							</div>
							<p v-if="titleError" class="help is-danger">
								{{ titleError }}
							</p>
						</div>

						<div class="field">
							<label class="label" for="merchantDescription">
								{{ $t('merchant.attributes.description') }}
							</label>
							<div class="control">
								<textarea
									id="merchantDescription"
									v-model="merchant.description"
									class="textarea"
									:placeholder="$t('merchant.edit.descriptionPlaceholder')"
									rows="4"
								/>
							</div>
						</div>

						<div class="field">
							<label class="label" for="merchantAddress">
								{{ $t('merchant.attributes.address') }}
							</label>
							<div class="control">
								<textarea
									id="merchantAddress"
									v-model="merchant.address"
									class="textarea"
									:placeholder="$t('merchant.edit.addressPlaceholder')"
									rows="3"
								/>
							</div>
						</div>

						<div class="columns">
							<div class="column">
								<div class="field">
									<label class="label" for="merchantPhone">
										{{ $t('merchant.attributes.phone') }}
									</label>
									<div class="control">
										<input
											id="merchantPhone"
											v-model="merchant.phone"
											class="input"
											:placeholder="$t('merchant.edit.phonePlaceholder')"
											type="tel"
											maxlength="20"
										>
									</div>
								</div>
							</div>
							<div class="column">
								<div class="field">
									<label class="label" for="merchantCity">
										{{ $t('merchant.attributes.city') }}
									</label>
									<div class="control">
										<input
											id="merchantCity"
											v-model="merchant.city"
											class="input"
											:placeholder="$t('merchant.edit.cityPlaceholder')"
											type="text"
											maxlength="50"
										>
									</div>
								</div>
							</div>
							<div class="column">
								<div class="field">
									<label class="label" for="merchantArea">
										{{ $t('merchant.attributes.area') }}
									</label>
									<div class="control">
										<input
											id="merchantArea"
											v-model="merchant.area"
											class="input"
											:placeholder="$t('merchant.edit.areaPlaceholder')"
											type="text"
											maxlength="50"
										>
									</div>
								</div>
							</div>
						</div>

						<div class="columns">
							<div class="column">
								<div class="field">
									<label class="label" for="merchantLng">
										{{ $t('merchant.attributes.lng') }}
									</label>
									<div class="control">
										<input
											id="merchantLng"
											v-model.number="merchant.lng"
											class="input"
											:placeholder="$t('merchant.edit.lngPlaceholder')"
											type="number"
											step="any"
											min="-180"
											max="180"
										>
									</div>
								</div>
							</div>
							<div class="column">
								<div class="field">
									<label class="label" for="merchantLat">
										{{ $t('merchant.attributes.lat') }}
									</label>
									<div class="control">
										<input
											id="merchantLat"
											v-model.number="merchant.lat"
											class="input"
											:placeholder="$t('merchant.edit.latPlaceholder')"
											type="number"
											step="any"
											min="-90"
											max="90"
										>
									</div>
								</div>
							</div>
							<div class="column is-narrow">
								<div class="field">
									<label class="label">&nbsp;</label>
									<div class="control">
										<button
											type="button"
											class="button is-info"
											:class="{ 'is-loading': geocoding }"
											:disabled="!merchant.address || geocoding"
											@click="geocodeAddress"
										>
											<icon icon="map-marker-alt" />
											{{ $t('merchant.geocode.button') }}
										</button>
									</div>
								</div>
							</div>
						</div>

						<div class="field">
							<div class="control">
								<label class="checkbox">
									<input
										v-model="merchant.isFavorite"
										type="checkbox"
									>
									{{ $t('merchant.attributes.favorite') }}
								</label>
							</div>
						</div>

						<div class="field is-grouped">
							<div class="control">
								<button
									type="submit"
									class="button is-success"
									:class="{ 'is-loading': loading }"
									:disabled="!merchant.title || loading"
								>
									<icon icon="save" />
									{{ $t('merchant.edit.save') }}
								</button>
							</div>
							<div class="control">
								<router-link
									:to="{ name: 'merchants.show', params: { merchantId: merchant.id } }"
									class="button"
								>
									{{ $t('misc.cancel') }}
								</router-link>
							</div>
							<div class="control">
								<button
									type="button"
									class="button is-danger"
									:class="{ 'is-loading': deleting }"
									:disabled="deleting"
									@click="deleteMerchant"
								>
									<icon icon="trash-alt" />
									{{ $t('merchant.edit.delete') }}
								</button>
							</div>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import {ref, reactive, computed, onMounted} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'

import MerchantModel from '@/models/merchant'
import MerchantService from '@/services/merchant'

import {success, error as showError} from '@/message'
import {useTitle} from '@/composables/useTitle'

const {t} = useI18n({useScope: 'global'})
const route = useRoute()
const router = useRouter()

const merchantService = new MerchantService()

const loading = ref(false)
const deleting = ref(false)
const geocoding = ref(false)
const merchant = reactive(new MerchantModel())

const merchantId = computed(() => parseInt(route.params.merchantId as string))

useTitle(() => merchant.title ? t('merchant.edit.titleWithName', {name: merchant.title}) : t('merchant.edit.title'))

const titleError = computed(() => {
	if (!merchant.title) return ''
	if (merchant.title.length > 100) {
		return t('merchant.edit.titleTooLong')
	}
	return ''
})

onMounted(async () => {
	await loadMerchant()
})

async function loadMerchant() {
	loading.value = true
	try {
		const loadedMerchant = await merchantService.get(merchantId.value)
		Object.assign(merchant, loadedMerchant)
	} catch (e) {
		console.error('Error loading merchant:', e)
		showError({message: t('merchant.edit.loadError')})
	} finally {
		loading.value = false
	}
}

async function updateMerchant() {
	if (titleError.value) return
	
	loading.value = true
	try {
		const updatedMerchant = await merchantService.update(merchant)
		Object.assign(merchant, updatedMerchant)
		success({message: t('merchant.edit.success')})
	} catch (e) {
		console.error('Error updating merchant:', e)
		showError({message: t('merchant.edit.updateError')})
	} finally {
		loading.value = false
	}
}

async function deleteMerchant() {
	if (!confirm(t('merchant.edit.deleteConfirm', {name: merchant.title}))) {
		return
	}
	
	deleting.value = true
	try {
		await merchantService.delete(merchant)
		success({message: t('merchant.edit.deleteSuccess')})
		await router.push({name: 'merchants.index'})
	} catch (e) {
		console.error('Error deleting merchant:', e)
		showError({message: t('merchant.edit.deleteError')})
	} finally {
		deleting.value = false
	}
}

async function geocodeAddress() {
	if (!merchant.address) return
	
	geocoding.value = true
	try {
		const result = await merchantService.geocode(merchant.address)
		if (result.longitude && result.latitude) {
			merchant.lng = result.longitude
			merchant.lat = result.latitude
			success({message: t('merchant.geocode.success')})
		}
	} catch (e) {
		console.error('Error geocoding address:', e)
		showError({message: t('merchant.geocode.error')})
	} finally {
		geocoding.value = false
	}
}
</script>

<style lang="scss" scoped>
.is-required {
	color: var(--danger);
}

.field {
	margin-bottom: 1.5rem;
}

.columns {
	margin-bottom: 1rem;
}

.card {
	margin: 1rem;
}
</style>
