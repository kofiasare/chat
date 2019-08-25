<template>
  <div id="login-container">
    <div v-if="websocket" id="login-wrapper">
      <el-form
        :model="ruleForm"
        :rules="rules"
        ref="ruleForm"
        @submit.prevent.native="submitForm('ruleForm')"
      >
        <el-form-item prop="chatname" label="Chat Name">
          <el-input
            @submit.prevent.native="handleLogin"
            v-model="ruleForm.chatname"
            placeholder="enter chat name"
            clearable
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            @click.native.prevent="submitForm('ruleForm')"
            :disabled="btnStatus"
            :loading="btnLoading"
          >Let's Chat</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-else>
      <div id="ws-missing">Your browser does not support WebSockets.</div>
    </div>
  </div>
</template>


<script>
export default {
  data() {
    return {
      websocket: true,
      btnLoading: false,
      ruleForm: {
        chatname: ""
      },
      rules: {
        chatname: [
          {
            required: true
          },

          {
            min: 4,
            max: 10
          }
        ]
      }
    };
  },

  methods: {
    submitForm(form) {
      this.btnLoading = true;
      this.$refs[form].validate(valid => {
        if (valid) {
          this.$router.push({
            path: "/chat",
            query: {
              u: this.ruleForm.chatname
            }
          });
          this.btnLoading = false;
        } else {
          this.btnLoading = false;
          return false;
        }
      });
    }
  },

  computed: {
    btnStatus() {
      return this.ruleForm.chatname === "";
    }
  },

  created() {
    window["WebSocket"] !== undefined
      ? (this.websocket = true)
      : (this.websocket = false);
  }
};
</script>


<style lang="scss">
#login-container {
  max-width: 400px;
  margin: 5rem auto;
  background-color: #fff;
}

#login-wrapper,
#ws-missing {
  padding: 2rem;
  box-shadow: 0 0 5px #ddd;
  border-radius: 5px;
}

#login-wrapper .el-button {
  width: 100%;
}
</style>
       