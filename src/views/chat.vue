<template>
  <div id="chat-container">
    <audio src="/sounds/tink.wav" />
    <el-container style="height: 100vh">
      <el-aside width="400px" style="background-color: #073042">
        <div>
          <p id="network-status">Currently Online</p>
          <div id="chat-users-container">
            <ul id="chat-users">
              <li v-for="user in users">{{user}}</li>
            </ul>
          </div>
        </div>
      </el-aside>

      <el-container>
        <el-main>
          <div id="chat-log"></div>
          <el-form id="chat-form" @submit.native.prevent="sendMesg">
            <el-row :gutter="20">
              <el-col :span="16">
                <el-form-item>
                  <el-input
                    type="textarea"
                    v-model="mesg"
                    placeholder="Enter your message here"
                    :autosize="true"
                  ></el-input>
                </el-form-item>
              </el-col>

              <el-col :span="8">
                <el-form-item>
                  <el-button type="primary" @click.prevent="sendMesg" :disabled="btnStatus">SEND</el-button>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  data() {
    return {
      ws: null,
      chatname: "",
      users: [],
      mesg: ""
    };
  },

  methods: {
    getUsers() {
      this.$http
        .get(`http://${document.location.host}/chat/api/users`)
        .then(res => (this.users = res.body))
        .catch(e => console.log(e));
    },

    appendLog(sender, mesg) {
      this.getUsers();
      const log = document.getElementById("chat-log");
      const audio = document.querySelector("audio");
      const scroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;

      const mesgDiv = document.createElement("div");
      const senderDiv = document.createElement("div");

      if (mesg) {
        senderDiv.className = "sender";
        mesgDiv.className = "mesg";

        senderDiv.innerHTML = sender;
        mesgDiv.innerHTML = mesg;

        log.appendChild(senderDiv);
        log.appendChild(mesgDiv);

        if (audio && sender !== this.chatname) {
          audio.currentTime = 0;
          audio.play();
        }

        if (scroll) log.scrollTop = log.scrollHeight - log.clientHeight;
      }
    },

    sendMesg() {
      this.ws.send(
        JSON.stringify({
          Sender: this.chatname,
          Text: this.mesg
        })
      );

      this.mesg = "";
    }
  },

  watch: {
    $route: {
      handler: function(route) {
        const query = route.query;
        if (query) {
          this.chatname = query.u;
        }
      },

      immediate: true
    }
  },

  created() {
    const wsURL = `ws://${document.location.host}/ws?u=${this.chatname}`;
    if (this.ws === null) {
      this.ws = new WebSocket(wsURL);
    }

    this.ws.onopen = e => this.getUsers();

    this.ws.onmessage = e => {
      const { sender, text } = JSON.parse(e.data);
      this.appendLog(sender, text);
    };

    this.ws.onerror = e => console.error(e);
    this.ws.onclose = () => this.appendLog("<b>Connection closed.</b>");
  },

  computed: {
    btnStatus() {
      return this.mesg === "";
    }
  }
};
</script>

<style lang="scss">
#chat-log {
  background: white;
  margin: 40px 20px 20px 20px;
  padding: 0.5em;
  height: 700px;
  overflow-y: auto;
}

#chat-form {
  margin: 20px;
}

#chat-container p {
  padding: 20px;
  text-align: center;
  font-weight: bolder;
}

#chat-container .el-aside {
  overflow: hidden;
}

#chat-container #chat-users {
  padding: 0;
  list-style: none;
}

#chat-container #chat-users li {
  padding: 10px 30px;
  margin: 15px;
  height: 27px;
  line-height: 27px;
  box-shadow: 0 0 5px #ddd;
}

#chat-users-container {
  margin: 10px;
  background-color: #f4f4f4;
  border-radius: 5px;
  height: 800px;
  overflow-y: auto;
  box-shadow: 0 0 5px #ddd;
}

.mesg {
  padding: 15px;
  background-color: azure;
  border-radius: 5px;
  box-shadow: 0 0 5px #ddd;
  margin: 9px 15px 25px 15px;
}

.sender {
  margin: 10px 20px;
  font-weight: bold;
}

#network-status {
  color: aliceblue;
}
</style>
   
