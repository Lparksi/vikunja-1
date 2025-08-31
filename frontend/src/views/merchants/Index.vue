<template>
	<div class="content-widescreen">
		<div class="merchant-management">
			<div class="header">
				<h1>{{ $t('merchant.title') }}</h1>
				<div class="actions">
					<BaseButton
						:to="{ name: 'merchants.create' }"
						variant="primary"
					>
						<Icon icon="plus" />
						{{ $t('merchant.create.title') }}
					</BaseButton>
					<BaseButton
						:to="{ name: 'merchants.import' }"
						variant="secondary"
					>
						<Icon icon="upload" />
						{{ $t('merchant.import.title') }}
					</BaseButton>
				</div>
			</div>

			<div class="filters">
				<div class="field">
					<div class="control has-icons-left">
						<input
							v-model="searchTerm"
							class="input"
							type="text"
							:placeholder="$t('merchant.search.placeholder')"
							@input="search"
						>
						<span class="icon is-left">
							<Icon icon="search" />
						</span>
					</div>
				</div>
				<div class="field">
					<div class="control">
						<div class="select">
							<select v-model="viewMode" @change="changeViewMode">
								<option value="list">{{ $t('merchant.view.list') }}</option>
								<option value="map">{{ $t('merchant.view.map') }}</option>
							</select>
						</div>
					</div>
				</div>
			</div>

			<div v-if="loading" class="loader-container">
				<div class="loader is-loading"></div>
			</div>

			<div v-else-if="viewMode === 'list'" class="merchant-list">
				<div
					v-for="merchant in merchants"
					:key="merchant.id"
					class="merchant-card"
				>
					<div class="merchant-info">
						<h3 class="merchant-title">{{ merchant.title }}</h3>
						<p class="merchant-address">{{ merchant.fullAddress }}</p>
						<div class="merchant-tags">
							<span
								v-for="tag in merchant.tags"
								:key="tag.id"
								class="tag"
								:style="{ backgroundColor: tag.hexColor }"
							>
								{{ tag.tagName }}
							</span>
						</div>
					</div>
					<div class="merchant-actions">
						<BaseButton
							:to="{ name: 'merchants.edit', params: { id: merchant.id } }"
							variant="secondary"
							size="small"
						>
							<Icon icon="edit" />
							{{ $t('misc.edit') }}
						</BaseButton>
						<BaseButton
							@click="deleteMerchant(merchant)"
							variant="danger"
							size="small"
						>
							<Icon icon="trash" />
							{{ $t('misc.delete') }}
						</BaseButton>
					</div>
				</div>

				<div v-if="merchants.length === 0" class="no-merchants">
					<p>{{ $t('merchant.empty') }}</p>
				</div>
			</div>

			<div v-else-if="viewMode === 'map'" class="merchant-map">
				<MerchantMap :merchants="merchants" />
			</div>

			<Pagination
				v-if="totalPages > 1"
				:total-pages="totalPages"
				:current-page="currentPage"
				@update:current-page="loadMerchants"
			/>
		</div>
	</div>
</template>

<script setup lang="ts">
import {ref, computed, onMounted} from 'vue'
import {useI18n} from 'vue-i18n'
import {useDebounceFn} from '@vueuse/core'

import BaseButton from '@/components/base/BaseButton.vue'
import Pagination from '@/components/misc/Pagination.vue'
import MerchantMap from '@/components/merchants/MerchantMap.vue'

import MerchantService from '@/services/merchant'
import MerchantModel from '@/models/merchant'
import type {IMerchant} from '@/modelTypes/IMerchant'

import {useTitle} from '@/composables/useTitle'
import {success, error} from '@/message'

const {t} = useI18n({useScope: 'global'})
useTitle(() => t('merchant.title'))

const merchantService = new MerchantService()

const merchants = ref<IMerchant[]>([])
const loading = computed(() => merchantService.loading)
const searchTerm = ref('')
const viewMode = ref('list')
const currentPage = ref(1)
const totalPages = computed(() => merchantService.totalPages)

const search = useDebounceFn(() => {
	currentPage.value = 1
	loadMerchants()
}, 500)

async function loadMerchants() {
	try {
		merchants.value = await merchantService.getAll(new MerchantModel(), {
			s: searchTerm.value,
		}, currentPage.value)
	} catch (e) {
		error({message: t('merchant.load.error')})
	}
}

function changeViewMode() {
	// View mode changed, could trigger map initialization
}

async function deleteMerchant(merchant: IMerchant) {
	if (!confirm(t('merchant.delete.confirm', {merchant: merchant.title}))) {
		return
	}

	try {
		await merchantService.delete(merchant)
		success({message: t('merchant.delete.success')})
		await loadMerchants()
	} catch (e) {
		error({message: t('merchant.delete.error')})
	}
}

onMounted(() => {
	loadMerchants()
})
</script>

<style lang="scss" scoped>
.merchant-management {
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;

		h1 {
			margin: 0;
		}

		.actions {
			display: flex;
			gap: 0.5rem;
		}
	}

	.filters {
		display: flex;
		gap: 1rem;
		margin-bottom: 1rem;
		align-items: center;

		.field {
			margin-bottom: 0;
		}
	}

	.loader-container {
		display: flex;
		justify-content: center;
		padding: 2rem;
	}

	.merchant-list {
		.merchant-card {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 1rem;
			border: 1px solid var(--grey-200);
			border-radius: 4px;
			margin-bottom: 0.5rem;
			background: var(--white);

			&:hover {
				box-shadow: var(--shadow-sm);
			}

			.merchant-info {
				flex: 1;

				.merchant-title {
					margin: 0 0 0.25rem 0;
					font-weight: 600;
				}

				.merchant-address {
					margin: 0 0 0.5rem 0;
					color: var(--grey-600);
					font-size: 0.9rem;
				}

				.merchant-tags {
					display: flex;
					gap: 0.25rem;
					flex-wrap: wrap;

					.tag {
						padding: 0.125rem 0.5rem;
						border-radius: 12px;
						font-size: 0.75rem;
						color: var(--white);
						background: var(--primary);
					}
				}
			}

			.merchant-actions {
				display: flex;
				gap: 0.5rem;
			}
		}

		.no-merchants {
			text-align: center;
			padding: 2rem;
			color: var(--grey-500);
		}
	}

	.merchant-map {
		height: 500px;
		border: 1px solid var(--grey-200);
		border-radius: 4px;
	}
}
</style>
