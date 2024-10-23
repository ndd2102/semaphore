<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="item != null && keys != null"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-text-field
        v-model="item.name"
        :label="$t('name')"
        :rules="[v => !!v || $t('name_required')]"
        required
        :disabled="formSaving"
    ></v-text-field>

    <v-text-field
        v-model.trim="item.url"
        :label="$t('urlOrPath')"
        :rules="[
          v => !!v || $t('repository_required'),
          v => getTypeOfUrl(v) != null || $t('incorrectUrl'),
        ]"
        required
        :disabled="formSaving"
        :hide-details="true"
    ></v-text-field>

    <v-text-field
      v-model.trim="item.bucket"
      :label="'bucket'"
      required
      :disabled="formSaving || type === 'local'"
    ></v-text-field>

    <v-text-field
      v-model.trim="item.access_key"
      :label="'access_key'"
      required
      :disabled="formSaving || type === 'local'"
    ></v-text-field>

    <v-text-field
      v-model.trim="item.secret_key"
      :label="'secret_key'"
      required
      type="password"
      :disabled="formSaving || type === 'local'"
    ></v-text-field>

  </v-form>
</template>
<script>
import axios from 'axios';
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],
  data() {
    return {
      helpDialog: null,
      helpKey: null,

      keys: null,
      inventoryTypes: [{
        id: 'static',
        name: 'Static',
      }, {
        id: 'static-yaml',
        name: 'Static YAML',
      }, {
        id: 'file',
        name: 'File',
      }],
    };
  },
  async created() {
    this.keys = (await axios({
      keys: 'get',
      url: `/api/project/${this.projectId}/keys`,
      responseType: 'json',
    })).data;
  },
  computed: {
    type() {
      return this.getTypeOfUrl(this.item.url);
    },
  },

  methods: {
    getTypeOfUrl(url) {
      if (url == null || url === '') {
        return null;
      }

      if (url.startsWith('/')) {
        return 'local';
      }

      const m = url.match(/^(\w+):\/\//);

      if (m == null) {
        return 'ssh';
      }

      if (!['git', 'file', 'ssh', 'http', 'https'].includes(m[1])) {
        return null;
      }

      return m[1];
    },

    setType(type) {
      let url;

      const m = this.item.url.match(/^\w+:\/\/(.*)$/);
      if (m != null) {
        url = m[1];
      } else {
        url = this.item.url;
      }

      if (type === 'local') {
        url = url.startsWith('/') ? url : `/${url}`;
      } else {
        url = `${type}://${url}`;
      }

      this.item.url = url;
    },

    showHelpDialog(key) {
      this.helpKey = key;
      this.helpDialog = true;
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/minioconfigs`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/minioconfigs/${this.itemId}`;
    },
  },
};
</script>
