<template>
  <div class="container">
    <h2>Conversor de Coordenadas + Seletor de Pasta</h2>

    <div class="input-section">
      <input v-model="input" />
      <button @click="converter">Adicionar</button>
    </div>

    <div class="folder-section">
      <button @click="selecionarPasta">Selecionar Pasta</button>
      <span v-if="pastaSelecionada" class="pasta-label">📁 {{ pastaSelecionada }}</span>
    </div>

    <table v-if="lista.length">
      <thead>
        <tr>
          <th>ID</th>
          <th>Coordenada Original</th>
          <th>Latitude</th>
          <th>Longitude</th>
          <th>Endereço</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in lista" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.coordenada }}</td>
          <td>{{ item.lat }}</td>
          <td>{{ item.long }}</td>
          <td>{{ item.endereco }}</td>
        </tr>
      </tbody>
    </table>

    <p v-else class="empty-message">Nenhuma coordenada adicionada ainda.</p>
    <button @click="enviarParaGerarKMZ">Gerar KMZ</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { SelectFolder } from '../../wailsjs/go/main/App.js' // Ajuste conforme sua estrutura
import { GerarKMZ } from '../../wailsjs/go/main/App.js'
const input = ref('')
const lista = ref([])
const pastaSelecionada = ref('')

const OPENCAGE_API_KEY = 'd9dcdd2fa85b4d58904d837b93209f5e'

function dmsParaDecimal(dms) {
  dms = dms.trim().toUpperCase()
  let sign = 1
  if (dms.includes('S') || dms.includes('W')) sign = -1

  dms = dms.replace(/[°'"]/g, ' ')
    .replace(/[NSEW]/g, '')
    .trim()

  const partes = dms.split(/\s+/)
  if (partes.length < 3) return null

  const graus = parseFloat(partes[0])
  const minutos = parseFloat(partes[1])
  const segundos = parseFloat(partes[2])
  if (isNaN(graus) || isNaN(minutos) || isNaN(segundos)) return null

  const decimal = graus + minutos / 60 + segundos / 3600
  return +(decimal * sign).toFixed(6)
}

function extrairCoordenadas(texto) {
  const regex = /(\d+°\d+'[\d.]+["]?[NS])\s+(\d+°\d+'[\d.]+["]?[EW])/i
  const match = texto.match(regex)
  if (!match || match.length < 3) return [null, null]
  return [match[1], match[2]]
}

async function buscarEndereco(lat, lon) {
  const url = `https://api.opencagedata.com/geocode/v1/json?q=${lat}+${lon}&key=${OPENCAGE_API_KEY}&language=pt`
  try {
    const res = await fetch(url)
    const data = await res.json()
    return data?.results?.[0]?.formatted || 'Endereço não encontrado'
  } catch (err) {
    console.error(err)
    return 'Erro na API'
  }
}

async function converter() {
  const [latStr, lonStr] = extrairCoordenadas(input.value)
  if (!latStr || !lonStr) {
    alert('Formato inválido. Exemplo: 23°34\'45.0"S 46°38\'20.0"W')
    return
  }

  const lat = dmsParaDecimal(latStr)
  const lon = dmsParaDecimal(lonStr)
  if (lat === null || lon === null) {
    alert('Erro ao converter coordenadas.')
    return
  }

  const endereco = await buscarEndereco(lat, lon)

  lista.value.push({
    id: lista.value.length + 1,
    coordenada: input.value,
    lat,
    long: lon,
    endereco
  })

  input.value = ''
}

async function selecionarPasta() {
  try {
    const path = await SelectFolder()
    if (path) pastaSelecionada.value = path
  } catch (err) {
    console.error('Erro ao selecionar pasta:', err)
  }
}
async function enviarParaGerarKMZ() {
  const payload = JSON.parse(JSON.stringify({
    pasta: pastaSelecionada.value,
    pontos: lista.value
  }))


  try {
    const res = await GerarKMZ(payload)
    alert("KMZ gerado com sucesso em: " + res)
  } catch (err) {
    console.error("Erro ao gerar KMZ:", err)
  }
}
</script>

<style scoped>
.container {
  max-width: 960px;
  margin: 2rem auto;
  padding: 1rem;
  font-family: Arial, sans-serif;
}

h2 {
  font-size: 24px;
  margin-bottom: 1rem;
  color: #333;
}

.input-section {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #bbb;
  border-radius: 4px;
  font-size: 14px;
}

button {
  padding: 10px 20px;
  font-size: 14px;
  background-color: #1976d2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #1565c0;
}

.folder-section {
  margin-bottom: 20px;
}

.pasta-label {
  margin-left: 10px;
  font-weight: bold;
  color: #333;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

th,
td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

thead {
  background-color: #f0f0f0;
}

tr:hover {
  background-color: #f9f9f9;
}

.empty-message {
  margin-top: 1rem;
  color: #777;
  font-style: italic;
}
</style>
