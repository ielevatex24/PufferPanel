<script setup>
import { ref, inject, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'

const { t } = useI18n()
const api = inject('api')
const templatesLoaded = ref(false)
const templates = ref([])
const firstEntry = ref(null)

onMounted(async () => {
  templates.value = await api.template.list()
  templatesLoaded.value = true
})

function setFirstEntry(ref) {
  if (!firstEntry.value) firstEntry.value = ref
}

function focusList() {
  firstEntry.value.$el.focus()
}
</script>

<template>
  <div class="templatelist">
    <h1 v-text="t('templates.Templates')" />
    <div v-hotkey="'l'" class="list" @hotkey="focusList()">
      <div v-for="template in templates" :key="template.name" class="list-item">
        <router-link :ref="setFirstEntry" :to="{ name: 'TemplateView', params: { id: template.name } }">
          <div class="template">
            <span class="title">{{template.display}}</span>
            <span class="subline">{{template.type}}</span>
          </div>
        </router-link>
      </div>
      <div v-if="!templatesLoaded" class="list-item">
        <loader small />
      </div>
      <div v-if="$api.auth.hasScope('templates.create')" class="list-item">
        <router-link v-hotkey="'c'" :to="{ name: 'TemplateCreate' }">
          <div class="createLink"><icon name="plus" />{{ t('templates.New') }}</div>
        </router-link>
      </div>
    </div>
  </div>
</template>
