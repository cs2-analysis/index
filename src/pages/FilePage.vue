<template>
  <q-page class="flex items-center column">
    <h3>{{ filename }}</h3>
    <q-table
      :title="route.params.path"
      :rows="rows"
      :columns="columns"
      row-key="name"
    >
      <template v-slot:body-cell-download="props">
        <q-td :props="props">
          <q-btn
            flat
            padding="6px"
            size="sm"
            icon="fas fa-download"
            :href="props.value"
            download="test"
          >
            <q-tooltip anchor="top middle" self="center middle">
              Original File
            </q-tooltip>
          </q-btn>
          <q-btn
            flat
            padding="6px"
            size="sm"
            icon="img:/img/ida.png"
          >
            <q-tooltip anchor="top middle" self="center middle">
              IDA Pro Database
            </q-tooltip>
          </q-btn>
          <q-btn
            flat
            padding="6px"
            size="sm"
            icon="img:/img/bindiff.ico"
          >
            <q-tooltip anchor="top middle" self="center middle">
              BinDiff Export
            </q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { computed, ref } from 'vue'
import windowsData from 'data/windows.json'
import linuxData from 'data/linux.json'
import { useRoute } from 'vue-router';
import { date, format } from 'quasar';

const columns = [
  {
    name: 'sha256',
    label: 'SHA256',
    field: 'sha256',
    align: 'left',
    format: val => `${val.slice(0, 8)}...`,
  },
  {
    name: 'manifest',
    label: 'Manifest',
    field: 'manifest',
    align: 'left',
  },
  {
    name: 'updated',
    label: 'Updated',
    field: 'updated',
    align: 'left',
    format: val => date.formatDate(val * 1000, 'YYYY-MM-DD HH:mm'),
  },
  {
    name: 'size',
    label: 'Size',
    field: 'size',
    align: 'left',
    format: val => format.humanStorageSize(val, 1)
      // add a space between the number and the unit
      .match(/^([\d.]+)\s*(\w+)$/).slice(1).join(' '),
  },
  {
    name: 'download',
    label: 'Download',
    align: 'center',
    format: (val, row) => `https://cdn.vaclive.party/cs2/${depot.value}/${row.sha256}/${row.filename}`,
  }
]

const route = useRoute()

const filename = computed(() => route.params.path.split('/').pop())

const rows = computed(() => {
  const data = route.params.platform === 'windows' ? windowsData : linuxData;
  return data[route.params.path];
})

const depot = computed(() => route.params.platform === 'windows' ? '2347771' : '2347773')

</script>
