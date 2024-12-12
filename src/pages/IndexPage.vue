<template>
  <q-page class="flex items-center column">

    <h3>CS:2 Binary Index</h3>

    <q-btn-toggle
      v-model="depot"
      :options="[
        {label: 'Windows', value: 'windows', icon: 'fab fa-windows'},
        {label: 'Linux', value: 'linux', icon: 'fab fa-linux'},
      ]"
      class="full-width q-mb-md"
      spread
      no-caps
      color="dark"
      text-color="white"
      padding="1rem"
    />

    <q-select
      filled
      model-value=""
      @update:model-value="opneFile"
      use-input
      hide-selected
      input-debounce="0"
      label="Search"
      :options="options"
      @filter="filterFn"
      class="full-width"
      hide-dropdown-icon
      fill-input
      emit-value
    >
      <template v-slot:append>
        <q-icon @click.stop.prevent name="fas fa-search" />
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="text-grey">
            No results
          </q-item-section>
        </q-item>
      </template>
    </q-select>

  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import windowsData from 'data/windows.json'
import linuxData from 'data/linux.json'
import { useRouter } from 'vue-router';

const getFilenames = data => Object.keys(data).map(key => ({
  value: key,
  label: key.split('/').pop(),
}));

const windowsFiles = getFilenames(windowsData)
const linuxFiles = getFilenames(linuxData)

const router = useRouter()

const depot = ref('windows')
const file = ref(null)
const options = ref([])

const filterFn = (val, update, abort) => {
  const optionsArr = depot.value === 'windows' ? windowsFiles : linuxFiles
  update(() => {
    if (val === '') {
      options.value = optionsArr
    }
    const needle = val.toLowerCase()
    options.value = optionsArr.filter(v => v.label.toLowerCase().indexOf(needle) > -1)
  })
}

const opneFile = (val) => {
  router.push(`/file/${depot.value}/${val}`);
}

</script>
