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
  props: ["chartData"],
  data (){
    return{
      labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
    }
  },
  methods: {
    test: function(){
      console.log(this.chartData);
      let data = this.chartData;
      this.renderChart({
        labels: data.Labels,
        datasets: [
          {
            backgroundColor: '#f87979',
            data: data.Value,
          }
        ]
      }, {responsive: true, maintainAspectRatio: false});
    }
  },
  mounted () {

    this.test();
  }

});

Vue.component("repo-info",{
  template: `<div><line-chart v-if=render v-bind:chartData="chartInfo"></line-chart><div class="pulsate" style="width:20px;height:20px;background-color:blue;left: 0; right: 0; margin: 0 auto;" v-if=vis></div></div>`,
  data (){
    return {
      repoData: "",
      repoDataOld: "blank",
      chartInfo: "blank",
      render: false,
      vis: false,

    }
  },
  methods: {
    update: function(){
      if(refresh){
        this.vis = true;
        refresh = false;
      }
      let el = this;
      this.repoData = repoData;
      if(this.repoData == this.repoDataOld){
        // console.log("woo
      }else{
        if(this.repoData !== ""){
          this.repoDataOld = this.repoData;
          $.get("/get/" + repoData,function(data){
            console.log(data);
            el.render = true;
            el.vis = false;
            el.chartInfo = data;
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
      refresh = true;
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
