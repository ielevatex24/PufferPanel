<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  idEditable: { type: Boolean, default: () => false },
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue', 'valid'])

const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
const nameInvalid = ref(false)
const displayInvalid = ref(false)
const typeInvalid = ref(false)

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function validate() {
  nameInvalid.value = false
  displayInvalid.value = false
  typeInvalid.value = false

  if (!template.value.name || template.value.name.trim() === '') {
    nameInvalid.value = true
  }

  if (!template.value.display || template.value.display.trim() === '') {
    displayInvalid.value = true
  }

  if (!template.value.type || template.value.type.trim() === '') {
    typeInvalid.value = true
  }

  emit('valid', !nameInvalid.value && !displayInvalid.value && !typeInvalid.value)
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u))
      template.value = u
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div>
    <text-field v-model="template.name" :disabled="!idEditable" :label="t('common.Name')" :hint="idEditable ? t('templates.description.Name') : undefined" :error="nameInvalid ? t('templates.errors.NameInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
    <text-field v-model="template.display" :label="t('templates.Display')" :hint="t('templates.description.Display')" :error="displayInvalid ? t('templates.errors.DisplayInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
    <text-field v-model="template.type" :label="t('templates.Type')" :hint="t('templates.description.Type')" :error="typeInvalid ? t('templates.errors.TypeInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
  </div>
</template>
