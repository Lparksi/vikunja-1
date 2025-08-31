<template>
	<div class="loader-container is-max-width-desktop" :class="{ 'is-loading': loading }">
		<div class="merchant-view">
			<!-- Header -->
			<div class="merchant-header">
				<div class="level">
					<div class="level-left">
						<div class="level-item">
							<div>
								<h1 class="title">
									{{ merchant.title }}
									<icon
										v-if="merchant.isFavorite"
										icon="star"
										class="favorite-icon"
									/>
								</h1>
								<p v-if="merchant.description" class="subtitle">
									{{ merchant.description }}
								</p>
							</div>
						</div>
					</div>
					<div class="level-right">
						<div class="level-item">
							<div class="buttons">
								<router-link
									:to="{ name: 'merchants.edit', params: { merchantId: merchant.id } }"
									class="button is-primary"
								>
									<icon icon="pen" />
									{{ $t('merchant.edit.title') }}
								</router-link>
								<button
									class="button"
									:class="{ 'is-success': !merchant.isFavorite, 'is-grey': merchant.isFavorite }"
									@click="toggleFavorite"
								>
									<icon :icon="merchant.isFavorite ? 'star' : ['far', 'star']" />
									{{ merchant.isFavorite ? $t('merchant.unfavorite') : $t('merchant.favorite') }}
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Content -->
			<div class="columns">
				<!-- Left Column - Details -->
				<div class="column is-two-thirds">
					<div class="card">
						<header class="card-header">
							<p class="card-header-title">
								{{ $t('merchant.details.title') }}
							</p>
						</header>
						<div class="card-content">
							<div class="content">
								<div class="merchant-details">
									<div v-if="merchant.address" class="detail-item">
										<strong>{{ $t('merchant.attributes.address') }}:</strong>
										<p>{{ merchant.address }}</p>
									</div>

									<div class="columns" v-if="merchant.city || merchant.area">
										<div v-if="merchant.city" class="column">
											<div class="detail-item">
												<strong>{{ $t('merchant.attributes.city') }}:</strong>
												<p>{{ merchant.city }}</p>
											</div>
										</div>
										<div v-if="merchant.area" class="column">
											<div class="detail-item">
												<strong>{{ $t('merchant.attributes.area') }}:</strong>
												<p>{{ merchant.area }}</p>
											</div>
										</div>
									</div>

									<div v-if="merchant.phone" class="detail-item">
										<strong>{{ $t('merchant.attributes.phone') }}:</strong>
										<p>
											<a :href="`tel:${merchant.phone}`">{{ merchant.phone }}</a>
										</p>
									</div>

									<div v-if="hasCoordinates" class="detail-item">
										<strong>{{ $t('merchant.attributes.coordinates') }}:</strong>
										<p>
											{{ $t('merchant.coordinates.format', { lng: merchant.lng, lat: merchant.lat }) }}
											<span v-if="merchant.geocodeDescription" class="tag is-small is-info ml-2">
												{{ merchant.geocodeDescription }}
											</span>
										</p>
									</div>

									<div v-if="merchant.tags && merchant.tags.length > 0" class="detail-item">
										<strong>{{ $t('merchant.attributes.tags') }}:</strong>
										<div class="tags">
											<span
												v-for="tag in merchant.tags"
												:key="tag.id"
												class="tag"
												:style="{ backgroundColor: tag.hexColor, color: getTextColor(tag.hexColor) }"
											>
												{{ tag.title }}
											</span>
										</div>
									</div>

									<div class="detail-item">
										<strong>{{ $t('merchant.attributes.created') }}:</strong>
										<p>{{ formatDate(merchant.created) }}</p>
									</div>

									<div v-if="merchant.updated !== merchant.created" class="detail-item">
										<strong>{{ $t('merchant.attributes.updated') }}:</strong>
										<p>{{ formatDate(merchant.updated) }}</p>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Geo Points -->
					<div v-if="merchant.geoPoints && merchant.geoPoints.length > 0" class="card mt-4">
						<header class="card-header">
							<p class="card-header-title">
								{{ $t('merchant.geoPoints.title') }}
							</p>
						</header>
						<div class="card-content">
							<div class="content">
								<div class="table-container">
									<table class="table is-fullwidth">
										<thead>
											<tr>
												<th>{{ $t('geoPoint.attributes.from') }}</th>
												<th>{{ $t('geoPoint.attributes.longitude') }}</th>
												<th>{{ $t('geoPoint.attributes.latitude') }}</th>
												<th>{{ $t('geoPoint.attributes.accuracy') }}</th>
												<th>{{ $t('geoPoint.attributes.created') }}</th>
											</tr>
										</thead>
										<tbody>
											<tr v-for="geoPoint in merchant.geoPoints" :key="geoPoint.id">
												<td>{{ geoPoint.from }}</td>
												<td>{{ geoPoint.longitude.toFixed(6) }}</td>
												<td>{{ geoPoint.latitude.toFixed(6) }}</td>
												<td>
													<span class="tag is-small" :class="getAccuracyClass(geoPoint.accuracy)">
														{{ geoPoint.accuracy }}%
													</span>
												</td>
												<td>{{ formatDate(geoPoint.created) }}</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- Right Column - Map/Actions -->
				<div class="column">
					<div class="card">
						<header class="card-header">
							<p class="card-header-title">
								{{ $t('merchant.actions.title') }}
							</p>
						</header>
						<div class="card-content">
							<div class="content">
								<div class="buttons is-vertical is-fullwidth">
									<button
										v-if="merchant.address"
										class="button is-info"
										:class="{ 'is-loading': geocoding }"
										:disabled="geocoding"
										@click="geocodeAddress"
									>
										<icon icon="map-marker-alt" />
										{{ $t('merchant.geocode.button') }}
									</button>

									<button
										v-if="hasCoordinates"
										class="button is-link"
										@click="openInMaps"
									>
										<icon icon="external-link-alt" />
										{{ $t('merchant.actions.openInMaps') }}
									</button>

									<router-link
										:to="{ name: 'merchants.edit', params: { merchantId: merchant.id } }"
										class="button is-primary"
									>
										<icon icon="pen" />
										{{ $t('merchant.edit.title') }}
									</router-link>

									<button
										class="button is-danger"
										@click="deleteMerchant"
									>
										<icon icon="trash-alt" />
										{{ $t('merchant.delete.title') }}
									</button>
								</div>
							</div>
						</div>
					</div>

					<!-- Map placeholder -->
					<div v-if="hasCoordinates" class="card mt-4">
						<header class="card-header">
							<p class="card-header-title">
								{{ $t('merchant.map.title') }}
							</p>
						</header>
						<div class="card-content">
							<div class="content">
								<div class="map-placeholder">
									<p class="has-text-centered">
										<icon icon="map" size="3x" />
									</p>
									<p class="has-text-centered mt-2">
										{{ $t('merchant.map.placeholder') }}
									</p>
									<p class="has-text-centered is-size-7 has-text-grey">
										{{ $t('merchant.coordinates.format', { lng: merchant.lng, lat: merchant.lat }) }}
									</p>
								</div>
							</div>
						</div>
					</div>
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
import {formatDateLong} from '@/helpers/time/formatDate'

const {t} = useI18n({useScope: 'global'})
const route = useRoute()
const router = useRouter()

const merchantService = new MerchantService()

const loading = ref(false)
const geocoding = ref(false)
const merchant = reactive(new MerchantModel())

const merchantId = computed(() => parseInt(route.params.merchantId as string))

useTitle(() => merchant.title || t('merchant.show.title'))

const hasCoordinates = computed(() => merchant.lng !== null && merchant.lat !== null)

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
		showError({message: t('merchant.show.loadError')})
	} finally {
		loading.value = false
	}
}

async function toggleFavorite() {
	try {
		merchant.isFavorite = !merchant.isFavorite
		await merchantService.update(merchant)
		success({
			message: merchant.isFavorite 
				? t('merchant.favorite.added') 
				: t('merchant.favorite.removed')
		})
	} catch (e) {
		merchant.isFavorite = !merchant.isFavorite // Revert on error
		console.error('Error toggling favorite:', e)
		showError({message: t('merchant.favorite.error')})
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
			await loadMerchant() // Reload to get updated geo points
		}
	} catch (e) {
		console.error('Error geocoding address:', e)
		showError({message: t('merchant.geocode.error')})
	} finally {
		geocoding.value = false
	}
}

async function deleteMerchant() {
	if (!confirm(t('merchant.delete.confirm', {name: merchant.title}))) {
		return
	}
	
	try {
		await merchantService.delete(merchant)
		success({message: t('merchant.delete.success')})
		await router.push({name: 'merchants.index'})
	} catch (e) {
		console.error('Error deleting merchant:', e)
		showError({message: t('merchant.delete.error')})
	}
}

function openInMaps() {
	if (!hasCoordinates.value) return
	
	const url = `https://www.google.com/maps?q=${merchant.lat},${merchant.lng}`
	window.open(url, '_blank')
}

function formatDate(date: string | Date) {
	return formatDateLong(new Date(date))
}

function getTextColor(hexColor: string): string {
	// Simple contrast calculation
	const r = parseInt(hexColor.slice(1, 3), 16)
	const g = parseInt(hexColor.slice(3, 5), 16)
	const b = parseInt(hexColor.slice(5, 7), 16)
	const brightness = (r * 299 + g * 587 + b * 114) / 1000
	return brightness > 128 ? '#000000' : '#ffffff'
}

function getAccuracyClass(accuracy: number): string {
	if (accuracy >= 80) return 'is-success'
	if (accuracy >= 60) return 'is-warning'
	return 'is-danger'
}
</script>

<style lang="scss" scoped>
.merchant-view {
	padding: 1rem;
}

.merchant-header {
	margin-bottom: 2rem;
	
	.title {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	.favorite-icon {
		color: var(--warning);
	}
}

.detail-item {
	margin-bottom: 1rem;
	
	strong {
		display: block;
		margin-bottom: 0.25rem;
		color: var(--text-light);
	}
	
	p {
		margin: 0;
	}
}

.tags {
	display: flex;
	flex-wrap: wrap;
	gap: 0.5rem;
}

.map-placeholder {
	min-height: 200px;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	background-color: var(--grey-100);
	border-radius: 4px;
	padding: 2rem;
}

.buttons.is-vertical {
	.button {
		margin-bottom: 0.5rem;
		
		&:last-child {
			margin-bottom: 0;
		}
	}
}

.table-container {
	overflow-x: auto;
}
</style>
