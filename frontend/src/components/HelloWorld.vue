<script setup>
import { reactive } from 'vue'
import { Greet } from '../../wailsjs/go/main/App'

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}

</script>

<template>
  <main>
  </main>
</template>

<style scoped></style>
