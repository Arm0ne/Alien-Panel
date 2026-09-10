import{$i as e,Cn as t,Fa as n,Gi as r,Hn as i,Oa as a,Qi as o,Qr as s,Sa as c,Si as l,Ta as u,Un as d,Vn as f,aa as p,bi as m,br as h,ea as g,fa as _,hi as v,ia as y,ir as b,ma as x,na as S,oa as C,qt as w,ra as T,ta as E,ua as D,wi as O,xa as k,xr as A,yi as j}from"./router-CGMPUmYf.js";import{r as M}from"./index-1M4Kl-DU.js";var N=j([j(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),m(`spin-container`,`
 position: relative;
 `,[m(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[f()])]),m(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),m(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[l(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),m(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),m(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[l(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=p({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},b.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),d),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=A(e),r=b(`Spin`,`-spin`,N,M,e,t),i=o(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{opacitySpinning:a,color:o,textColor:s}=i;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?v(t):i[O(`size`,t)],"--n-color":o,"--n-text-color":s}}),c=n?h(`spin`,o(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),i,e):void 0,l=s(e,[`spinning`,`show`]),u=a(!1);return k(t=>{let n;if(l.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{u.value=!0},r),t(()=>{clearTimeout(n)});return}}u.value=l.value}),{mergedClsPrefix:t,active:u,mergedStrokeWidth:o(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return P[typeof n==`number`?`medium`:n]}),cssVars:n?void 0:i,themeClass:c?.themeClass,onRender:c?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:a}=this,o=t.icon&&this.rotate,s=(a||t.description)&&C(`div`,{class:`${n}-spin-description`},a||t.description?.call(t)),c=t.icon?C(`div`,{class:[`${n}-spin-body`,this.themeClass]},C(`div`,{class:[`${n}-spin`,o&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),s):C(`div`,{class:[`${n}-spin-body`,this.themeClass]},C(i,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),s);return(e=this.onRender)==null||e.call(this),t.default?C(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},C(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),C(r,{name:`fade-in-transition`},{default:()=>this.active?c:null})):c}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(t,n){return x(),S(`svg`,I,[...n[0]||=[e(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=u({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=p({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(r){let i=r,s=a(Date.now()),l,u=o(()=>{if(!i.dataAt)return null;let e=new Date(i.dataAt).getTime();return Number.isFinite(e)?e:null}),d=o(()=>{if(u.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=s.value-u.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(i.delayedAfterMinutes,1)*60*1e3,n=Math.max(i.expiredAfterMinutes,i.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),f=o(()=>u.value===null?``:new Date(u.value).toLocaleString(`zh-CN`,{hour12:!1}));return _(()=>{i.compact||(l=setInterval(()=>{s.value=Date.now()},6e4))}),D(()=>{l&&clearInterval(l)}),(i,a)=>{let o=t,s=w;return x(),g(s,null,{trigger:c(()=>[y(o,{size:`small`,type:d.value.type},{default:c(()=>[T(n(d.value.label)+` `,1),!r.compact&&f.value?(x(),S(`span`,z,`· `+n(f.value),1)):E(``,!0)]),_:1},8,[`type`])]),default:c(()=>[e(`div`,null,n(d.value.detail),1),f.value?(x(),S(`div`,B,`数据时间：`+n(f.value),1)):E(``,!0)]),_:1})}}});export{R as n,F as r,V as t};