<template>
	<div class="merchant-map-container">
		<div class="map-placeholder">
			<div class="map-content">
				<div class="map-icon">
					<icon icon="map" size="3x" />
				</div>
				<h3 class="map-title">{{ $t('merchant.map.title') }}</h3>
				<p class="map-description">{{ $t('merchant.map.description') }}</p>
				
				<div v-if="merchants.length > 0" class="merchant-stats">
					<div class="stat">
						<span class="stat-number">{{ merchants.length }}</span>
						<span class="stat-label">{{ $t('merchant.map.totalMerchants') }}</span>
					</div>
					<div class="stat">
						<span class="stat-number">{{ merchantsWithCoordinates }}</span>
						<span class="stat-label">{{ $t('merchant.map.withCoordinates') }}</span>
					</div>
				</div>

				<div v-if="merchantsWithCoordinates > 0" class="merchant-list">
					<h4>{{ $t('merchant.map.merchantsWithLocation') }}</h4>
					<div class="merchant-items">
						<div
							v-for="merchant in merchantsWithLocation"
							:key="merchant.id"
							class="merchant-item"
						>
							<div class="merchant-info">
								<strong>{{ merchant.title }}</strong>
								<p class="merchant-address">{{ merchant.address }}</p>
								<div class="merchant-coordinates">
									<small>
										{{ $t('merchant.coordinates.format', { 
											lng: merchant.lng?.toFixed(6), 
											lat: merchant.lat?.toFixed(6) 
										}) }}
									</small>
								</div>
							</div>
							<div class="merchant-actions">
								<button
									class="button is-small is-info"
									@click="openInMaps(merchant)"
								>
									<icon icon="external-link-alt" />
									{{ $t('merchant.map.openInMaps') }}
								</button>
							</div>
						</div>
					</div>
				</div>

				<div v-else-if="merchants.length > 0" class="no-coordinates">
					<p>{{ $t('merchant.map.noCoordinates') }}</p>
					<p class="help-text">{{ $t('merchant.map.addCoordinatesHelp') }}</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {IMerchant} from '@/modelTypes/IMerchant'

const {t} = useI18n({useScope: 'global'})

interface Props {
	merchants: IMerchant[]
}

const props = defineProps<Props>()

const merchantsWithLocation = computed(() => {
	return props.merchants.filter(merchant => 
		merchant.lng !== null && 
		merchant.lat !== null &&
		merchant.lng !== undefined && 
		merchant.lat !== undefined
	)
})

const merchantsWithCoordinates = computed(() => merchantsWithLocation.value.length)

function openInMaps(merchant: IMerchant) {
	if (merchant.lng && merchant.lat) {
		const url = `https://www.google.com/maps?q=${merchant.lat},${merchant.lng}`
		window.open(url, '_blank')
	}
}
</script>

<style lang="scss" scoped>
.merchant-map-container {
	height: 100%;
	min-height: 400px;
	border: 1px solid var(--grey-200);
	border-radius: 8px;
	overflow: hidden;
}

.map-placeholder {
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, var(--grey-50) 0%, var(--grey-100) 100%);
	padding: 2rem;
}

.map-content {
	text-align: center;
	max-width: 600px;
	width: 100%;
}

.map-icon {
	color: var(--grey-400);
	margin-bottom: 1rem;
}

.map-title {
	margin: 0 0 0.5rem 0;
	color: var(--grey-700);
	font-size: 1.5rem;
	font-weight: 600;
}

.map-description {
	margin: 0 0 2rem 0;
	color: var(--grey-600);
	font-size: 1rem;
	line-height: 1.5;
}

.merchant-stats {
	display: flex;
	justify-content: center;
	gap: 2rem;
	margin-bottom: 2rem;
	
	.stat {
		text-align: center;
		
		.stat-number {
			display: block;
			font-size: 2rem;
			font-weight: 700;
			color: var(--primary);
			line-height: 1;
		}
		
		.stat-label {
			display: block;
			font-size: 0.875rem;
			color: var(--grey-600);
			margin-top: 0.25rem;
		}
	}
}

.merchant-list {
	text-align: left;
	
	h4 {
		margin: 0 0 1rem 0;
		color: var(--grey-700);
		font-size: 1.125rem;
		font-weight: 600;
	}
}

.merchant-items {
	max-height: 300px;
	overflow-y: auto;
	border: 1px solid var(--grey-200);
	border-radius: 6px;
	background: var(--white);
}

.merchant-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 1rem;
	border-bottom: 1px solid var(--grey-100);
	
	&:last-child {
		border-bottom: none;
	}
	
	&:hover {
		background: var(--grey-50);
	}
}

.merchant-info {
	flex: 1;
	
	strong {
		display: block;
		margin-bottom: 0.25rem;
		color: var(--grey-800);
	}
	
	.merchant-address {
		margin: 0 0 0.25rem 0;
		color: var(--grey-600);
		font-size: 0.875rem;
	}
	
	.merchant-coordinates {
		small {
			color: var(--grey-500);
			font-family: monospace;
		}
	}
}

.merchant-actions {
	margin-left: 1rem;
}

.no-coordinates {
	padding: 2rem;
	text-align: center;
	
	p {
		margin: 0 0 0.5rem 0;
		color: var(--grey-600);
	}
	
	.help-text {
		font-size: 0.875rem;
		color: var(--grey-500);
	}
}

@media (max-width: 768px) {
	.merchant-stats {
		flex-direction: column;
		gap: 1rem;
	}
	
	.merchant-item {
		flex-direction: column;
		align-items: flex-start;
		gap: 1rem;
		
		.merchant-actions {
			margin-left: 0;
			align-self: stretch;
		}
	}
}
</style>
