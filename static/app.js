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
  extends: VueChartJs.Bar,
  props: ["chartData"],
  data (){
    return{
      labels: [],
    }
  },
  methods: {
    test: function(){
      console.log(this.chartData);
      let data = this.chartData;
      text = []
      vals = []
      for (var i = 0; i < data.length; i ++){
        text.push(data[i].Label);

        vals.push(data[i].Value);

      }
      this.renderChart({
        labels: text,
        datasets: [
          {
            label: "",
            backgroundColor: '#f87979',
            data: vals,
          }
        ]
      }, {responsive: true, maintainAspectRatio: false, scales:{xAxes:[{barPercentage: 1}]}, legend:{display: false}});
    }
  },
  mounted () {

    this.test();
  }

});

Vue.component("repo-card",{
  template: `<div><a v-bind:href=data.url<div class="card">
  <img class="avatar" v-bind:src=data.owner.avatar_url></img>
  <a class="repoName">{{data.full_name}}</a>
  <p>Project Status: Who knows?</p>
  </div></div>`,
  props: ['data'],
  data (){
    return {

    }
  }
})

Vue.component("repo-info",{
  template: `<div>
  <div class="pulsate" style="width:20px;height:20px;background-color:blue;left: 0; right: 0; margin: 0 auto;" v-if=vis></div>
  <repo-card v-if="render && repoInfo.full_name != ''" v-bind:data=repoInfo ></repo-card>
  </div>`,
  data (){
    return {
      repoData: "",
      repoDataOld: "blank",
      repoInfo: "blank",
      render: false,
      vis: false,

    }
  },
  methods: {
    update: function(){
      if(refresh){
        this.vis = true;
        refresh = false;
        this.render = false;
      }
      let el = this;
      this.repoData = repoData;
      if(this.repoData == this.repoDataOld){
        // console.log("woo
      }else{
        if(this.repoData !== ""){
          this.repoDataOld = this.repoData;
          $.get("/get/" + repoData,function(data){
            if(data == null){

            }else{
              el.render = true;
              el.vis = false;
              el.repoInfo = data;

            }
          });

        }
      }
    }
  },
  mounted(){
    setInterval(this.update, 100);
  }
});

var timeout = null;
Vue.component("get-repo",{
  template: `<div v-bind:style=textStyle>
              <input v-bind:style=inputStyle @keyup=update v-bind:placeholder=repos v-model=textboxText>
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
      let el = this;
      clearTimeout(timeout);
      el.inputStyle.borderColor = "#e74c3c";
      timeout = setTimeout(function () {
        if(el.textboxText.indexOf("/") != -1){
          refresh = true;
          el.inputStyle.borderColor = "#1abc9c";
          repoData = el.textboxText;
        }
      }, 500);
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
