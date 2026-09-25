import{Ci as e,Dn as t,Ea as n,Ei as r,Gn as i,Kn as a,Ma as o,Oi as s,Ta as c,Tr as l,Yi as u,Zt as d,aa as f,ca as p,cr as m,ga as h,ia as g,ka as _,la as v,ma as y,na as b,ni as x,oa as S,qn as C,ra as w,sa as T,ua as E,va as D,wi as O,wr as k,yi as A,za as j}from"./router-BPt_aGUh.js";import{r as M}from"./index-BmdhIupJ.js";var N=e([e(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),O(`spin-container`,`
 position: relative;
 `,[O(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[i()])]),O(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),O(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[r(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),O(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),O(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[r(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=v({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},m.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),C),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=l(e),r=m(`Spin`,`-spin`,N,M,e,t),i=b(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{opacitySpinning:a,color:o,textColor:c}=i;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?A(t):i[s(`size`,t)],"--n-color":o,"--n-text-color":c}}),a=n?k(`spin`,b(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),i,e):void 0,u=x(e,[`spinning`,`show`]),d=o(!1);return c(t=>{let n;if(u.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{d.value=!0},r),t(()=>{clearTimeout(n)});return}}d.value=u.value}),{mergedClsPrefix:t,active:d,mergedStrokeWidth:b(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return P[typeof n==`number`?`medium`:n]}),cssVars:n?void 0:i,themeClass:a?.themeClass,onRender:a?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,o=(r||t.description)&&E(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),s=t.icon?E(`div`,{class:[`${n}-spin-body`,this.themeClass]},E(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),o):E(`div`,{class:[`${n}-spin-body`,this.themeClass]},E(a,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),o);return(e=this.onRender)==null||e.call(this),t.default?E(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},E(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),E(u,{name:`fade-in-transition`},{default:()=>this.active?s:null})):s}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(e,t){return D(),S(`svg`,I,[...t[0]||=[w(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=_({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=v({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let r=e,i=o(Date.now()),a,s=b(()=>{if(!r.dataAt)return null;let e=new Date(r.dataAt).getTime();return Number.isFinite(e)?e:null}),c=b(()=>{if(s.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-s.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(r.delayedAfterMinutes,1)*60*1e3,n=Math.max(r.expiredAfterMinutes,r.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),l=b(()=>s.value===null?``:new Date(s.value).toLocaleString(`zh-CN`,{hour12:!1}));return h(()=>{r.compact||(a=setInterval(()=>{i.value=Date.now()},6e4))}),y(()=>{a&&clearInterval(a)}),(r,i)=>{let a=t,o=d;return D(),g(o,null,{trigger:n(()=>[p(a,{size:`small`,type:c.value.type},{default:n(()=>[T(j(c.value.label)+` `,1),!e.compact&&l.value?(D(),S(`span`,z,`· `+j(l.value),1)):f(``,!0)]),_:1},8,[`type`])]),default:n(()=>[w(`div`,null,j(c.value.detail),1),l.value?(D(),S(`div`,B,`数据时间：`+j(l.value),1)):f(``,!0)]),_:1})}}});export{R as n,F as r,V as t};