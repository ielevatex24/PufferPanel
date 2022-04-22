<script setup>
import { ref, inject, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import AnsiUp from 'ansi_up'

const ansiup = new AnsiUp()
ansiup.ansi_to_html('\u001b[m')

const { t } = useI18n()
const config = inject('config')
const name = config.branding.name

const command = ref('')
const console = ref(null)
let lastMessageTime = 0
let lastIncompleteLine = null

const props = defineProps({
  server: { type: Object, required: true }
})

let unbindEvent = null
let task = null
onMounted(() => {
  unbindEvent = props.server.on('console', onConsole)

  task = props.server.startTask(() => {
    if (props.server.needsPolling()) {
      props.server.replayConsole(lastMessageTime)
    }
  }, 5000)

  props.server.replayConsole()
})

onUnmounted(() => {
  if (unbindEvent) unbindEvent()
  if (task) props.server.stopTask(task)
  clearConsole()
})

function handleCarriageReturn(line) {
  if (line.indexOf('\r') !== -1) {
    const parts = line.split('\r')
    let result = parts.shift()
    parts.map(part => {
      result = part + result.substring(part.length)
    })
    return result
  }

  return line
}

function markDaemon(line) {
  if (line.trim().indexOf('[DAEMON]') === 0) {
    return `<span class="daemon-marker" data-name="${name}"></span>` + line.substring(8)
  }

  return line
}

function handleLine(line) {
  // escaping after carriage return to not mess with char counts
  return markDaemon(handleCarriageReturn(line))
}

function onConsole(event) {
  if ('epoch' in event) {
    lastMessageTime = event.epoch
  } else {
    lastMessageTime = Math.floor(Date.now() / 1000)
  }

  let newLines = (Array.isArray(event.logs) ? event.logs.join('') : event.logs).replaceAll('\r\n', '\n')
  const endOnNewline = newLines.endsWith('\n')
  newLines = newLines.split('\n')

  if (endOnNewline) {
    // if ending on a newline, do not render an empty last line
    newLines.pop()
  }

  let last = null
  newLines.map(line => {
    line = ansiup.ansi_to_html(line)
    if (lastIncompleteLine) {
      line = lastIncompleteLine.line + line
    }

    if (lastIncompleteLine) {
      lastIncompleteLine.el.innerHTML = handleLine(line)
      last = { el: lastIncompleteLine.el, line }
    } else {
      const el = document.createElement('div')
      el.innerHTML = handleLine(line)
      console.value.appendChild(el)
      last = { el, line }
    }
  })

  if (!endOnNewline) {
    // not ending on a newline means last line is incomplete
    // therefore remember it to complete later
    lastIncompleteLine = last
  } else {
    lastIncompleteLine = null
  }
}

function clearConsole() {
  if (console.value) console.value.replaceChildren([])
}

function sendCommand() {
  props.server.sendCommand(command.value)
  command.value = ''
}
</script>

<template>
  <div>
    <h2 v-text="t('servers.Console')" />
    <icon v-hotkey="'c c'" name="clear-console" @click="clearConsole()" />
    <div dir="ltr" class="console-wrapper">
      <div ref="console" class="console" />
    </div>
    <div dir="ltr" class="command">
      <text-field v-model="command" :label="t('servers.Command')" @keyup.enter="sendCommand()" />
      <icon name="send" @click="sendCommand()" />
    </div>
  </div>
</template>
