<template>
	<div class="loader-container is-max-width-desktop" :class="{ 'is-loading': loading }">
		<div class="card">
			<header class="card-header">
				<p class="card-header-title">
					{{ $t('merchant.import.title') }}
				</p>
			</header>
			<div class="card-content">
				<div class="content">
					<!-- Import Instructions -->
					<div class="notification is-info">
						<h4 class="title is-5">{{ $t('merchant.import.instructions.title') }}</h4>
						<div class="content">
							<ol>
								<li>{{ $t('merchant.import.instructions.step1') }}</li>
								<li>{{ $t('merchant.import.instructions.step2') }}</li>
								<li>{{ $t('merchant.import.instructions.step3') }}</li>
								<li>{{ $t('merchant.import.instructions.step4') }}</li>
							</ol>
						</div>
					</div>

					<!-- Template Download -->
					<div class="field">
						<label class="label">{{ $t('merchant.import.template.title') }}</label>
						<div class="control">
							<button
								class="button is-info"
								@click="downloadTemplate"
							>
								<icon icon="download" />
								{{ $t('merchant.import.template.download') }}
							</button>
						</div>
						<p class="help">{{ $t('merchant.import.template.description') }}</p>
					</div>

					<!-- File Upload -->
					<div class="field">
						<label class="label">{{ $t('merchant.import.file.title') }}</label>
						<div class="file has-name" :class="{ 'is-danger': fileError }">
							<label class="file-label">
								<input
									ref="fileInput"
									class="file-input"
									type="file"
									accept=".csv"
									@change="handleFileSelect"
								>
								<span class="file-cta">
									<span class="file-icon">
										<icon icon="upload" />
									</span>
									<span class="file-label">
										{{ $t('merchant.import.file.choose') }}
									</span>
								</span>
								<span v-if="selectedFile" class="file-name">
									{{ selectedFile.name }}
								</span>
							</label>
						</div>
						<p v-if="fileError" class="help is-danger">
							{{ fileError }}
						</p>
						<p v-else class="help">
							{{ $t('merchant.import.file.description') }}
						</p>
					</div>

					<!-- Import Options -->
					<div class="field">
						<label class="label">{{ $t('merchant.import.options.title') }}</label>
						<div class="control">
							<label class="checkbox">
								<input
									v-model="autoGeocode"
									type="checkbox"
								>
								{{ $t('merchant.import.options.autoGeocode') }}
							</label>
						</div>
						<p class="help">{{ $t('merchant.import.options.autoGeocodeDescription') }}</p>
					</div>

					<!-- Import Button -->
					<div class="field">
						<div class="control">
							<button
								class="button is-success is-large"
								:class="{ 'is-loading': importing }"
								:disabled="!selectedFile || importing || fileError"
								@click="importMerchants"
							>
								<icon icon="upload" />
								{{ $t('merchant.import.start') }}
							</button>
						</div>
					</div>

					<!-- Import Results -->
					<div v-if="importResult" class="notification" :class="getResultClass()">
						<h4 class="title is-5">{{ $t('merchant.import.result.title') }}</h4>
						<div class="content">
							<div class="columns">
								<div class="column">
									<div class="has-text-centered">
										<p class="heading">{{ $t('merchant.import.result.total') }}</p>
										<p class="title">{{ importResult.totalRows }}</p>
									</div>
								</div>
								<div class="column">
									<div class="has-text-centered">
										<p class="heading">{{ $t('merchant.import.result.success') }}</p>
										<p class="title has-text-success">{{ importResult.successCount }}</p>
									</div>
								</div>
								<div class="column">
									<div class="has-text-centered">
										<p class="heading">{{ $t('merchant.import.result.errors') }}</p>
										<p class="title has-text-danger">{{ importResult.errorCount }}</p>
									</div>
								</div>
							</div>

							<!-- Error Details -->
							<div v-if="importResult.errors && importResult.errors.length > 0" class="mt-4">
								<h5 class="title is-6">{{ $t('merchant.import.result.errorDetails') }}</h5>
								<div class="table-container">
									<table class="table is-fullwidth is-striped">
										<thead>
											<tr>
												<th>{{ $t('merchant.import.result.row') }}</th>
												<th>{{ $t('merchant.import.result.column') }}</th>
												<th>{{ $t('merchant.import.result.message') }}</th>
												<th>{{ $t('merchant.import.result.value') }}</th>
											</tr>
										</thead>
										<tbody>
											<tr v-for="error in importResult.errors" :key="`${error.row}-${error.column}`">
												<td>{{ error.row }}</td>
												<td>{{ error.column }}</td>
												<td>{{ error.message }}</td>
												<td class="is-family-monospace">{{ error.value || '-' }}</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>

							<!-- Success Actions -->
							<div v-if="importResult.successCount > 0" class="mt-4">
								<div class="buttons">
									<router-link
										:to="{ name: 'merchants.index' }"
										class="button is-primary"
									>
										<icon icon="list" />
										{{ $t('merchant.import.result.viewMerchants') }}
									</router-link>
									<button
										class="button is-info"
										@click="resetImport"
									>
										<icon icon="plus" />
										{{ $t('merchant.import.result.importMore') }}
									</button>
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
import {ref, computed} from 'vue'
import {useI18n} from 'vue-i18n'

import MerchantService from '@/services/merchant'

import {success, error as showError} from '@/message'
import {useTitle} from '@/composables/useTitle'

const {t} = useI18n({useScope: 'global'})
useTitle(() => t('merchant.import.title'))

const merchantService = new MerchantService()

const loading = ref(false)
const importing = ref(false)
const selectedFile = ref<File | null>(null)
const autoGeocode = ref(true)
const importResult = ref<any>(null)
const fileInput = ref<HTMLInputElement>()

const fileError = computed(() => {
	if (!selectedFile.value) return ''
	
	if (!selectedFile.value.name.toLowerCase().endsWith('.csv')) {
		return t('merchant.import.file.invalidType')
	}
	
	if (selectedFile.value.size > 10 * 1024 * 1024) { // 10MB limit
		return t('merchant.import.file.tooLarge')
	}
	
	return ''
})

function handleFileSelect(event: Event) {
	const target = event.target as HTMLInputElement
	if (target.files && target.files.length > 0) {
		selectedFile.value = target.files[0]
		importResult.value = null // Clear previous results
	}
}

function downloadTemplate() {
	const template = `商户名称,地址,电话,城市,区域,描述,标签,经度,纬度
示例商户,北京市朝阳区建国门外大街1号,010-12345678,北京,朝阳区,这是一个示例商户,"餐饮,服务",116.397128,39.916527`

	const blob = new Blob([template], { type: 'text/csv;charset=utf-8;' })
	const link = document.createElement('a')
	const url = URL.createObjectURL(blob)
	
	link.setAttribute('href', url)
	link.setAttribute('download', 'merchant_import_template.csv')
	link.style.visibility = 'hidden'
	
	document.body.appendChild(link)
	link.click()
	document.body.removeChild(link)
	
	success({message: t('merchant.import.template.downloaded')})
}

async function importMerchants() {
	if (!selectedFile.value || fileError.value) return
	
	importing.value = true
	importResult.value = null
	
	try {
		const result = await merchantService.importFromCSV(selectedFile.value, autoGeocode.value)
		importResult.value = result
		
		if (result.successCount > 0) {
			success({
				message: t('merchant.import.success', {
					count: result.successCount,
					total: result.totalRows
				})
			})
		}
		
		if (result.errorCount > 0) {
			showError({
				message: t('merchant.import.partialError', {
					errors: result.errorCount,
					total: result.totalRows
				})
			})
		}
	} catch (e) {
		console.error('Error importing merchants:', e)
		showError({message: t('merchant.import.error')})
	} finally {
		importing.value = false
	}
}

function resetImport() {
	selectedFile.value = null
	importResult.value = null
	autoGeocode.value = true
	
	if (fileInput.value) {
		fileInput.value.value = ''
	}
}

function getResultClass() {
	if (!importResult.value) return ''
	
	if (importResult.value.errorCount === 0) {
		return 'is-success'
	} else if (importResult.value.successCount > 0) {
		return 'is-warning'
	} else {
		return 'is-danger'
	}
}
</script>

<style lang="scss" scoped>
.card {
	margin: 1rem;
}

.file.has-name .file-name {
	max-width: 300px;
	overflow: hidden;
	text-overflow: ellipsis;
}

.table-container {
	max-height: 400px;
	overflow-y: auto;
}

.notification .content {
	margin-bottom: 0;
}

.buttons {
	margin-top: 1rem;
}

.heading {
	font-size: 0.75rem;
	font-weight: 600;
	text-transform: uppercase;
	color: var(--text-light);
}
</style>
