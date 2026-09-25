import{Ci as e,Cr as t,Di as n,En as r,Gn as i,Ji as a,Kn as o,Oa as s,Ra as c,Si as l,Ta as u,Ti as d,Wn as f,Xt as p,_a as m,aa as h,ca as g,ha as _,ia as v,ja as y,la as b,na as x,oa as S,pa as C,ra as w,sa as T,sr as E,ta as D,ti as O,vi as k,wa as A,wr as j}from"./router-CZSFYymG.js";import{r as M}from"./index-C1zbMUbF.js";var N=l([l(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),e(`spin-container`,`
 position: relative;
 `,[e(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[f()])]),e(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),e(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[d(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),e(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),e(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[d(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=g({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},E.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),o),slots:Object,setup(e){let{mergedClsPrefixRef:r,inlineThemeDisabled:i}=j(e),a=E(`Spin`,`-spin`,N,M,e,r),o=D(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:r},self:i}=a.value,{opacitySpinning:o,color:s,textColor:c}=i;return{"--n-bezier":r,"--n-opacity-spinning":o,"--n-size":typeof t==`number`?k(t):i[n(`size`,t)],"--n-color":s,"--n-text-color":c}}),s=i?t(`spin`,D(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),o,e):void 0,c=O(e,[`spinning`,`show`]),l=y(!1);return A(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{l.value=!0},r),t(()=>{clearTimeout(n)});return}}l.value=c.value}),{mergedClsPrefix:r,active:l,mergedStrokeWidth:D(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return P[typeof n==`number`?`medium`:n]}),cssVars:i?void 0:o,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,o=t.icon&&this.rotate,s=(r||t.description)&&b(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),c=t.icon?b(`div`,{class:[`${n}-spin-body`,this.themeClass]},b(`div`,{class:[`${n}-spin`,o&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),s):b(`div`,{class:[`${n}-spin-body`,this.themeClass]},b(i,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),s);return(e=this.onRender)==null||e.call(this),t.default?b(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},b(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),b(a,{name:`fade-in-transition`},{default:()=>this.active?c:null})):c}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(e,t){return m(),h(`svg`,I,[...t[0]||=[x(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=s({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=g({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let t=e,n=y(Date.now()),i,a=D(()=>{if(!t.dataAt)return null;let e=new Date(t.dataAt).getTime();return Number.isFinite(e)?e:null}),o=D(()=>{if(a.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=n.value-a.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let r=Math.max(t.delayedAfterMinutes,1)*60*1e3,i=Math.max(t.expiredAfterMinutes,t.delayedAfterMinutes)*60*1e3;return e>i?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(i/6e4)} 分钟未同步`}:e>r?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),s=D(()=>a.value===null?``:new Date(a.value).toLocaleString(`zh-CN`,{hour12:!1}));return _(()=>{t.compact||(i=setInterval(()=>{n.value=Date.now()},6e4))}),C(()=>{i&&clearInterval(i)}),(t,n)=>{let i=r,a=p;return m(),w(a,null,{trigger:u(()=>[T(i,{size:`small`,type:o.value.type},{default:u(()=>[S(c(o.value.label)+` `,1),!e.compact&&s.value?(m(),h(`span`,z,`· `+c(s.value),1)):v(``,!0)]),_:1},8,[`type`])]),default:u(()=>[x(`div`,null,c(o.value.detail),1),s.value?(m(),h(`div`,B,`数据时间：`+c(s.value),1)):v(``,!0)]),_:1})}}});export{R as n,F as r,V as t};