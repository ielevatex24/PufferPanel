<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import markdown from '@/utils/markdown.js'

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')
const emit = defineEmits(['selected'])
const available = ref([])
const importable = ref([])
const showing = ref(false)
const currentTemplate = ref({})

async function load() {
  const known = await api.template.list()
  available.value = known
  const knownRemote = await api.template.listImportable()
  importable.value = knownRemote.filter(i => !known.find(t => t.name === i))
}

onMounted(async () => {
  load()
})

async function importAndShow(template) {
  if (await api.template.import(template)) {
    load()
    show(template)
  } else {
    toast.error(t('errors.ImportFailed'))
  }
}

async function show(t) {
  currentTemplate.value = await api.template.get(t)
  showing.value = true
}

function choice(confirm) {
  showing.value = false
  if (confirm) emit('selected', currentTemplate.value)
}
</script>

<template>
  <div class="select-template">
    <h2 v-text="t('servers.SelectTemplate')" />
    <div v-for="template in available" :key="template.name" class="template" @click="show(template.name)">
      <span v-text="template.display" />
    </div>
    <h3 v-text="t('templates.Import')" />
    <div class="warning" v-text="t('templates.import.CommunityWarning')" />
    <div v-for="template in importable" :key="template" class="template importable" @click="importAndShow(template)">
      <span v-text="template" />
    </div>

    <overlay v-model="showing" :title="currentTemplate.display" closable>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-if="currentTemplate.readme" dir="ltr" class="readme" v-html="markdown(currentTemplate.readme)" />
      <h2 v-else v-text="t('servers.ConfirmTemplateChoice')" />
      <div class="actions">
        <btn color="error" @click="choice(false)"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="choice(true)"><icon name="check" />{{ t('servers.SelectThisTemplate') }}</btn>
      </div>
    </overlay>
  </div>
</template>
