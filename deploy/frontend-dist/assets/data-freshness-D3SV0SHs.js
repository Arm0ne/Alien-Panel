import{$i as e,Bn as t,Hi as n,Ma as r,Qi as i,Rn as a,Sa as o,Ta as s,Wt as c,Xi as l,Yi as u,Yr as d,Zi as f,_i as p,_r as m,bn as h,da as g,ea as _,fi as v,gi as y,la as b,na as x,ra as S,sa as C,ta as w,tr as T,va as E,vr as D,xi as O,ya as k,yi as A,zn as j}from"./router-Dj5ayp1M.js";import{r as M}from"./index-DFs6bn7n.js";var N=x({name:`Backward`,render(){return S(`svg`,{viewBox:`0 0 20 20`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},S(`path`,{d:`M12.2674 15.793C11.9675 16.0787 11.4927 16.0672 11.2071 15.7673L6.20572 10.5168C5.9298 10.2271 5.9298 9.7719 6.20572 9.48223L11.2071 4.23177C11.4927 3.93184 11.9675 3.92031 12.2674 4.206C12.5673 4.49169 12.5789 4.96642 12.2932 5.26634L7.78458 9.99952L12.2932 14.7327C12.5789 15.0326 12.5673 15.5074 12.2674 15.793Z`,fill:`currentColor`}))}}),P=x({name:`FastBackward`,render(){return S(`svg`,{viewBox:`0 0 20 20`,version:`1.1`,xmlns:`http://www.w3.org/2000/svg`},S(`g`,{stroke:`none`,"stroke-width":`1`,fill:`none`,"fill-rule":`evenodd`},S(`g`,{fill:`currentColor`,"fill-rule":`nonzero`},S(`path`,{d:`M8.73171,16.7949 C9.03264,17.0795 9.50733,17.0663 9.79196,16.7654 C10.0766,16.4644 10.0634,15.9897 9.76243,15.7051 L4.52339,10.75 L17.2471,10.75 C17.6613,10.75 17.9971,10.4142 17.9971,10 C17.9971,9.58579 17.6613,9.25 17.2471,9.25 L4.52112,9.25 L9.76243,4.29275 C10.0634,4.00812 10.0766,3.53343 9.79196,3.2325 C9.50733,2.93156 9.03264,2.91834 8.73171,3.20297 L2.31449,9.27241 C2.14819,9.4297 2.04819,9.62981 2.01448,9.8386 C2.00308,9.89058 1.99707,9.94459 1.99707,10 C1.99707,10.0576 2.00356,10.1137 2.01585,10.1675 C2.05084,10.3733 2.15039,10.5702 2.31449,10.7254 L8.73171,16.7949 Z`}))))}}),F=x({name:`FastForward`,render(){return S(`svg`,{viewBox:`0 0 20 20`,version:`1.1`,xmlns:`http://www.w3.org/2000/svg`},S(`g`,{stroke:`none`,"stroke-width":`1`,fill:`none`,"fill-rule":`evenodd`},S(`g`,{fill:`currentColor`,"fill-rule":`nonzero`},S(`path`,{d:`M11.2654,3.20511 C10.9644,2.92049 10.4897,2.93371 10.2051,3.23464 C9.92049,3.53558 9.93371,4.01027 10.2346,4.29489 L15.4737,9.25 L2.75,9.25 C2.33579,9.25 2,9.58579 2,10.0000012 C2,10.4142 2.33579,10.75 2.75,10.75 L15.476,10.75 L10.2346,15.7073 C9.93371,15.9919 9.92049,16.4666 10.2051,16.7675 C10.4897,17.0684 10.9644,17.0817 11.2654,16.797 L17.6826,10.7276 C17.8489,10.5703 17.9489,10.3702 17.9826,10.1614 C17.994,10.1094 18,10.0554 18,10.0000012 C18,9.94241 17.9935,9.88633 17.9812,9.83246 C17.9462,9.62667 17.8467,9.42976 17.6826,9.27455 L11.2654,3.20511 Z`}))))}}),I=x({name:`Forward`,render(){return S(`svg`,{viewBox:`0 0 20 20`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},S(`path`,{d:`M7.73271 4.20694C8.03263 3.92125 8.50737 3.93279 8.79306 4.23271L13.7944 9.48318C14.0703 9.77285 14.0703 10.2281 13.7944 10.5178L8.79306 15.7682C8.50737 16.0681 8.03263 16.0797 7.73271 15.794C7.43279 15.5083 7.42125 15.0336 7.70694 14.7336L12.2155 10.0005L7.70694 5.26729C7.42125 4.96737 7.43279 4.49264 7.73271 4.20694Z`,fill:`currentColor`}))}}),L=y([y(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),p(`spin-container`,`
 position: relative;
 `,[p(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[a()])]),p(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),p(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[A(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),p(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),p(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[A(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),R={small:20,medium:18,large:16},z=x({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},T.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),t),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=D(e),r=T(`Spin`,`-spin`,L,M,e,t),i=u(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{opacitySpinning:a,color:o,textColor:s}=i;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?v(t):i[O(`size`,t)],"--n-color":o,"--n-text-color":s}}),a=n?m(`spin`,u(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),i,e):void 0,o=d(e,[`spinning`,`show`]),c=s(!1);return E(t=>{let n;if(o.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{c.value=!0},r),t(()=>{clearTimeout(n)});return}}c.value=o.value}),{mergedClsPrefix:t,active:c,mergedStrokeWidth:u(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return R[typeof n==`number`?`medium`:n]}),cssVars:n?void 0:i,themeClass:a?.themeClass,onRender:a?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:r,description:i}=this,a=t.icon&&this.rotate,o=(i||t.description)&&S(`div`,{class:`${r}-spin-description`},i||t.description?.call(t)),s=t.icon?S(`div`,{class:[`${r}-spin-body`,this.themeClass]},S(`div`,{class:[`${r}-spin`,a&&`${r}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),o):S(`div`,{class:[`${r}-spin-body`,this.themeClass]},S(j,{clsPrefix:r,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${r}-spin`}),o);return(e=this.onRender)==null||e.call(this),t.default?S(`div`,{class:[`${r}-spin-container`,this.themeClass],style:this.cssVars},S(`div`,{class:[`${r}-spin-content`,this.active&&`${r}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),S(n,{name:`fade-in-transition`},{default:()=>this.active?s:null})):s}}),B={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function V(t,n){return g(),e(`svg`,B,[...n[0]||=[l(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var H=o({name:`mdi-refresh`,render:V}),U={key:0,class:`ml-4px`},W={key:0,class:`mt-2px text-12px opacity-80`},G=x({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let n=t,a=s(Date.now()),o,d=u(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),p=u(()=>{if(d.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=a.value-d.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,r=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>r?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),m=u(()=>d.value===null?``:new Date(d.value).toLocaleString(`zh-CN`,{hour12:!1}));return b(()=>{n.compact||(o=setInterval(()=>{a.value=Date.now()},6e4))}),C(()=>{o&&clearInterval(o)}),(n,a)=>{let o=h,s=c;return g(),f(s,null,{trigger:k(()=>[w(o,{size:`small`,type:p.value.type},{default:k(()=>[_(r(p.value.label)+` `,1),!t.compact&&m.value?(g(),e(`span`,U,`· `+r(m.value),1)):i(``,!0)]),_:1},8,[`type`])]),default:k(()=>[l(`div`,null,r(p.value.detail),1),m.value?(g(),e(`div`,W,`数据时间：`+r(m.value),1)):i(``,!0)]),_:1})}}});export{F as a,I as i,H as n,P as o,z as r,N as s,G as t};