<template>
  <v-text-field
    :id="id"
    ref="underlying"
    :value="value"
    :outlined="lookValue === 'outlined'"
    :solo="lookValue === 'solo'"
    :filled="lookValue === 'filled'"
    :flat="flat"
    :autofocus="autofocus"
    :autocomplete="autocomplete"
    :dense="dense"
    :disabled="disabled"
    :label="label"
    :hide-details="hideDetails && !showHint"
    :persistent-hint="showHint"
    :hint="hintValue"
    :error-messages="errorMessages"
    :prepend-inner-icon="icon"
    :append-icon="endIcon"
    :append-outer-icon="iconBehind"
    :placeholder="placeholder"
    :required="required"
    :name="name"
    :type="type"
    @input="$emit('input', $event)"
    v-on="listeners"
  >
    <slot
      v-for="(_, slotName) in $slots"
      :slot="slotName"
      :name="slotName"
    />
  </v-text-field>
</template>

<script>
const allowedLooks = ['outlined', 'solo', 'solo-flat', 'filled', 'material']

export default {
  props: {
    autofocus: { type: Boolean, default: () => false },
    autocomplete: { type: String, default: () => undefined },
    dense: { type: Boolean, default: () => false },
    disabled: { type: Boolean, default: () => false },
    endIcon: { type: String, default: () => undefined },
    errorMessages: { type: String, default: () => undefined },
    hint: { type: String, default: () => undefined },
    icon: { type: String, default: () => undefined },
    iconBehind: { type: String, default: () => undefined },
    id: { type: String, default: () => undefined },
    label: { type: String, default: () => undefined },
    name: { type: String, default: () => undefined },
    placeholder: { type: String, default: () => undefined },
    required: { type: Boolean, default: () => false },
    look: { type: String, validator: val => { return allowedLooks.indexOf(val) !== -1 }, default: () => undefined },
    type: { type: String, default: () => 'text' },
    value: { type: undefined, default: () => '', required: true }
  },
  computed: {
    listeners () {
      const { input, ...listeners } = this.$listeners
      return listeners
    },
    hideDetails () {
      const hasError =
        this.errorMessages !== undefined && this.errorMessages !== ''
      return !hasError
    },
    showHint () {
      const hasHint = this.hint !== undefined && this.hint !== ''
      const hasMsgSlot =
        this.$slots.message !== undefined &&
        this.$slots.message !== '' &&
        this.$slots.message !== []
      return hasHint || hasMsgSlot
    },
    lookValue () {
      const defaulted = this.look ? this.look : this.$vuetify.theme.options.inputStyle
      return defaulted.split('-')[0]
    },
    flat () {
      const defaulted = this.look ? this.look : this.$vuetify.theme.options.inputStyle
      return defaulted.split('-').indexOf('flat') !== -1
    },
    hintValue () {
      // set hint to '_' if only the slot has content to force vuetify to
      // display the hint without needing to double define it everywhere
      return this.hint ? this.hint : this.$slots.message ? '_' : undefined
    }
  },
  methods: {
    focus () {
      this.$refs.underlying.focus()
    }
  }
}
</script>
