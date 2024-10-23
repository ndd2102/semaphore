<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items != null && keys != null">
    <EditDialog
      v-model="editDialog"
      :save-button-text="itemId === 'new' ? $t('create') : $t('save')"
      :title="`${itemId === 'new' ? $t('nnew') : $t('edit')} minioconfigs`"
      @save="loadItems()"
      :max-width="450"
    >
      <template v-slot:form="{ onSave, onError, needSave, needReset }">
        <MinioForm
          :project-id="projectId"
          :item-id="itemId"
          @save="onSave"
          @error="onError"
          :need-save="needSave"
          :need-reset="needReset"
        />
      </template>
    </EditDialog>

    <ObjectRefsDialog
      object-title="minioconfigs"
      :object-refs="itemRefs"
      :project-id="projectId"
      v-model="itemRefsDialog"
    />

    <YesNoDialog
      :title="'minioconfigs'"
      :text="$t('askDeleteRepo')"
      v-model="deleteItemDialog"
      @yes="deleteItem(itemId)"
    />

    <v-toolbar flat >
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>minioconfigs</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        @click="editItem('new')"
        v-if="can(USER_PERMISSIONS.manageProjectResources)"
      >minioconfigs</v-btn>
    </v-toolbar>

    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
      class="mt-4"
      :items-per-page="Number.MAX_VALUE"
    >
      <template v-slot:item.git_url="{ item }">
        {{ item.git_url }}
        <code v-if="!item.git_url.startsWith('/')">{{ item.git_branch }}</code>
      </template>

      <template v-slot:item.ssh_key_id="{ item }">
        {{ keys.find((k) => k.id === item.ssh_key_id).name }}
      </template>

      <template v-slot:item.actions="{ item }">
        <div style="white-space: nowrap">
          <v-btn
            icon
            class="mr-1"
            @click="askDeleteItem(item.id)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            icon
            class="mr-1"
            @click="editItem(item.id)"
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
        </div>
      </template>
    </v-data-table>
  </div>

</template>
<script>
import ItemListPageBase from '@/components/ItemListPageBase';
import MinioForm from '@/components/MinioForm.vue';
import axios from 'axios';

export default {
  mixins: [ItemListPageBase],
  components: { MinioForm },
  data() {
    return {
      keys: null,
    };
  },

  async created() {
    this.keys = (await axios({
      method: 'get',
      url: `/api/project/${this.projectId}/keys`,
      responseType: 'json',
    })).data;
  },

  mounted() {
    console.log(this.getItemsUrl());
  },
  methods: {
    getHeaders() {
      return [{
        text: this.$i18n.t('name'),
        value: 'name',
        width: '25%',
      },
      {
        text: 'URL',
        value: 'url',
        width: '50%',
      },
      {
        text: 'bucket',
        value: 'bucket',
        width: '25%',
      },
      {
        text: this.$i18n.t('actions'),
        value: 'actions',
        sortable: false,
      }];
    },
    getItemsUrl() {
      return `/api/project/${this.projectId}/minioconfigs`;
    },
    getSingleItemUrl() {
      return `/api/project/${this.projectId}/minioconfigs/${this.itemId}`;
    },
    getEventName() {
      return 'i-minioconfigs';
    },
  },
};
</script>
