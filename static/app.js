var repoData = "";
var refresh = false;

Vue.component("header-text",{
  template: `<div v-bind:style=textStyle><h1 v-bind:style=headerstyle>BurnDown</h1>
  <p v-bind:style="subtext">Use burn down charts as a project sanity check</p></div>`,
  data (){
    return{
      textStyle: {textAlign:"center",position:"absolute",width:"auto",left:"0",right:"0",top:"100px"},
      headerstyle: {fontWeight: "lighter",fontSize:"40px"},
      subtext: {color:"#444",fontSize:"20px"}
    }
  }
});

Vue.component('line-chart', {
  extends: VueChartJs.Line,
  props: ["data"],
  data (){
    return{
      labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
    }
  },
  methods: {
    test: function(){
    }
  },
  mounted () {
    this.renderChart({
      labels: this.labels,
      datasets: [
        {
          label: 'Burndown',
          backgroundColor: '#f87979',
          data: [30, 20, 30, 20, 10, 5, 1]
        }
      ]
    }, {responsive: true, maintainAspectRatio: false});

  }
});

Vue.component("repo-info",{
  template: `<div><line-chart v-if=render></line-chart></div>`,
  data (){
    return {
      repoData: "",
      repoDataOld: "asdf",
      other: "adf",
      render: false,
    }
  },
  methods: {
    update: function(){
      let el = this;
      this.repoData = repoData;
      if(this.repoData == this.repoDataOld){

      }else{
        if(this.repoData !== ""){
          this.render = true;
          this.repoDataOld = this.repoData;
          $.get("https://api.github.com/repos/" + repoData,function(data){
            el.other = data;
          });

        }
      }
    }
  },
  mounted(){
    setInterval(this.update, 100);
  }
});

Vue.component("get-repo",{
  template: `<div v-bind:style=textStyle>
              <input v-bind:style=inputStyle @keyup.enter=update v-bind:placeholder=repos v-model=textboxText>
              </input>

              </div>`,
  data (){
    return {
      textStyle: {textAlign:"center",position:"absolute",width:"auto",left:"0",right:"0",top:"250px"},
      inputStyle: {height: "50px",textAlign:"center", borderWidth:"1px",borderStyle:"solid",borderRadius:"5px",borderColor:"#1abc9c",width:"50%",fontSize:"30px",padding:"30px"},
      repos: "wtg/shuttletracker",
      textboxText:""
    }
  },
  methods: {
    update: function(){
      repoData = this.textboxText;
    }
  },

});


Vue.component("titlebar",{
  template: `<div v-bind:style=titleStyle><p v-bind:style=paragraphStyle>{{titleText}}</p></div>`,
  data (){
    return{
      titleStyle: {position:"absolute",backgroundColor:"#1abc9c",height:"50px",width:"auto",top:"1",left:"0",right:"0",color:"white"},
      paragraphStyle: {float: "left",height:"50px",lineHeight:"50px",verticalAlign:"center",paddingLeft:"30px",margin:"0"},
      titleText: "BurnDown 🔥"
    }
  }

});

var App = new Vue({
  el: '#app-vue',
  data: {
  },


});
