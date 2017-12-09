var repoData = "";
var refresh = false;

Vue.component("header-text",{
  template: `<div v-bind:style=textStyle><h1 v-bind:style=headerstyle>BurnDown</h1>
  <p v-bind:style="subtext">Use charts as a project sanity check - </p>

  <a href="https://github.com/jlyon1/burndown" v-bind:style="smaller">Source</a> - <a href="https://github.com/jlyon1/burndown" v-bind:style="smaller">About</a></div>`,
  data (){
    return{
      textStyle: {textAlign:"center",position:"absolute",width:"auto",left:"0",right:"0",top:"100px"},
      headerstyle: {fontWeight: "lighter",fontSize:"40px"},
      subtext: {color:"#444",fontSize:"20px"},
      smaller: {color:"#3498db",fontSize:"15px"}

    }
  }
});

Vue.component('bar-chart', {
  extends: VueChartJs.Line,
  data (){
    return{
      chartData: [],
    }
  },
  methods: {
    render: function(){
      let data = this.chartData;
      text = []
      vals = []
      for(let i = 0; i <data.length; i ++){
        text.push(data[i].Label)
        //console.log(data[i])
        if(i == 0){
          vals.push(data[i].Value)
        }else{
          vals.push(vals[i-1] + data[i].Value)
        }
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
      }, {responsive: true, maintainAspectRatio: false, scales:{xAxes:[{barPercentage: 1}]}, elements: { point: { radius: 0 } }, legend:{display: false}});
    }
  },
  mounted () {
    let el = this;
    $.get("/bar/" + repoData,function(data){
      el.chartData = data;
      el.render()
    }).fail(function(){
    });
  }

});

Vue.component("chart-card",{
  template: `<div>
  <div style="margin: 20px;" class="box">
  <span class="repoName">Burndown Chart:</span>
  <bar-chart></bar-chart>
  </div>
  </div>`,
  props: ['data'],
  data (){
    return {

    }
  }
})

Vue.component("repo-card",{
  template: `<div><a style="text-decoration:none; color:black;" v-bind:href=data.html_url>
  <div style="margin: 20px;" class="box">
  <span class="repoName">{{data.full_name}}</span>
  <img class="image is-64x64 media-right" v-bind:src=data.owner.avatar_url style="float:right;"></img>
  <p>Project Status: {{staleness.Text}}</p>
  <p stype="color:#aaa;">Loaded: {{data.Issues.length}} issues, {{data.Commits.length}} commits, and {{data.Pulls.length}} Pull Requests</p>
  </div>
  </a>
  </div>`,
  props: ['data'],
  data (){
    return {
      staleness: 0
    }
  },
  mounted(){
    let el = this;
    $.get("/stale/" + repoData, function(data){
      el.staleness = data;
    })
  }
})

Vue.component("single-issue",{
  props: ['issue','max'],
  template: `<div style="margin:10px;background-color:white;height:30px;width:100%;border-bottom-style:solid;border-width:1px;border-color:#eee;">
  <div style="width:90%;overflow:hide;float:left;"><a :href=issue.link>{{issue.label}}</a></div>
  <div style="width:10%;float:right;"><progress class="progress" v-bind:class="{'is-danger': (percent >= .75),'is-warning': (percent >= .5 && percent < .75),'is-success': (percent <.5)}" :value=this.percent :max=1>{{percent}}</progress></div>
  </div>`,
  data (){
    return {
      percent: this.issue.val/this.max,
    }
  }
})

Vue.component("issue-card",{
  props: ['data'],
  template: `<div><div style="height: auto; margin: 20px;overflow:hide;" class="box">
  <div style="font-size: 18px;margin-bottom: 10px;" class="repoName">Issues: {{data.Issues.length}} <span> <span style="color:#aaa;">-</span> Open Issues: {{issueData.Open}}</span><span> <span style="color:#aaa;">-</span> Average open time: {{(issueData.AvgDuration/(3600*24)).toFixed(2)}} days</span></div>
  Open:
  <single-issue v-for="val in open" :issue=val :max=issueData.MaxDuration></single-issue>
  Closed:

  <single-issue v-for="val in closed" :issue=val :max=issueData.MaxDuration></single-issue>

  </div>
  </div>`,
  data (){
    return {
      issueData: {},
      open: [],
      closed: []
    }
  },
  mounted (){
    let el = this;
    let obj = []
    let cls = []
    $.get("/issue/" + repoData,function(data){
      el.issueData = data;
      for(let i =0; i < data.Data[0].Points.length; i ++){
        if(data.Data[0].Points[i].Value < (.1*el.issueData.MaxDuration)){
          obj.push({link: data.Data[0].Points[i].Link, label: data.Data[0].Points[i].Label,val: (.1*el.issueData.MaxDuration)});
        }else{
          obj.push({link: data.Data[0].Points[i].Link, label: data.Data[0].Points[i].Label,val: data.Data[0].Points[i].Value});
        }
      }
      for(let i =0; i < data.Data[1].Points.length; i ++){
        if(data.Data[1].Points[i].Value < (.1*el.issueData.MaxDuration)){
          cls.push({link: data.Data[1].Points[i].Link, label: data.Data[1].Points[i].Label,val: (.1*el.issueData.MaxDuration)});
        }else{
          cls.push({link: data.Data[1].Points[i].Link, label: data.Data[1].Points[i].Label,val: data.Data[1].Points[i].Value});
        }
      }
      //console.log(cls);
    }).fail(function(){
    });
    this.open = obj;
    this.closed = cls;

    //console.log(this.open)
  }
})


Vue.component("repo-info",{
  template: `<div>
  <div class="pulsate" style="width:20px;height:20px;background-color:blue;left: 0; right: 0; margin: 0 auto;" v-if=vis></div>
  <repo-card v-if="render && repoInfo.full_name != ''" v-bind:data=repoInfo ></repo-card>
  <chart-card v-if="render && repoInfo.full_name != ''"></chart-card>
  <issue-card v-if="render && repoInfo.full_name != ''" v-bind:data=repoInfo></issue-card>

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
      repoData = window.location.pathname.substr(1,window.location.pathname.length);

      if(refresh){
        this.vis = true;
        refresh = false;
        this.repoData = ""
        this.render = false;
      }
      let el = this;
      this.repoData = repoData;
      if(this.repoData == this.repoDataOld && this.render){

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
          }).fail(function(){

          });;
        }
      }
    }
  },
  mounted(){
    //console.log(window.location.pathname);
    setInterval(this.update, 100);
  }
});

var timeout = null;
Vue.component("get-repo",{
  template: `<div v-bind:style=textStyle>
              <input class="input" style="width:50%;text-align:center;" @keyup=update v-bind:placeholder=repos v-model=textboxText>
              </input>

              </div>`,
  data (){
    return {
      textStyle: {textAlign:"center",position:"absolute",width:"auto",left:"0",right:"0",top:"250px"},
      inputStyle: {height: "50px",textAlign:"center", borderWidth:"1px",borderStyle:"solid",borderRadius:"5px",borderColor:"#1abc9c",width:"50%",fontSize:"30px",padding:"30px"},
      repos: "wtg/shuttletracker",
      textboxText:"wtg/shuttletracker"
    }
  },
  methods: {
    update: function(){
      let el = this;
      clearTimeout(timeout);
      el.inputStyle.borderColor = "#e74c3c";
      timeout = setTimeout(function () {
        if(el.textboxText.indexOf("/") != -1){
          $.get("/valid/" + el.textboxText,function(data){
            refresh = true;
            if(data == "true"){
              history.pushState({content:document.documentElement.innerHTML}, null, "/" + el.textboxText);
              el.inputStyle.borderColor = "#1abc9c";
              repoData = el.textboxText;
            }
          }).fail(function(){

          });;
        }
      }, 500);
    }

  },
  mounted(){
    this.textboxText = window.location.pathname.substr(1,window.location.pathname.length);
    this.update();

  }

});


Vue.component("titlebar",{
  template: `<div v-bind:style=titleStyle><p v-bind:style=paragraphStyle>{{titleText}}</p></div>`,
  data (){
    return{
      titleStyle: {position:"absolute",backgroundColor:"#1abc9c",height:"50px",width:"auto",top:"0.1",left:"0",right:"0",color:"white"},
      paragraphStyle: {float: "left",height:"50px",lineHeight:"50px",verticalAlign:"center",paddingLeft:"30px",margin:"0"},
      titleText: "🔥⬇️"
    }
  }

});

var App = new Vue({
  el: '#app-vue',
  data: {
  },


});
