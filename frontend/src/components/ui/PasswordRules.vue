<script setup>
import { computed } from 'vue'
import CheckIcon from './CheckIcon.vue'
import { checkPassword } from '../../utils/validators.js'

const props = defineProps({ password: { type: String, default: '' } })
const rules = computed(() => checkPassword(props.password))
const score = computed(() => rules.value.filter((r) => r.ok).length)
</script>

<template>
  <div class="pw-rules">
    <div class="pw-meter" :data-score="score" aria-hidden="true"><span v-for="n in 4" :key="n" :class="{ on: n <= score }" /></div>
    <ul>
      <li v-for="r in rules" :key="r.id" :class="{ ok: r.ok }">
        <CheckIcon :on="r.ok" />{{ r.label }}
      </li>
    </ul>
  </div>
</template>
