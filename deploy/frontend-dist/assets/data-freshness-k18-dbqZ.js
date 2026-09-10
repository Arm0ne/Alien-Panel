import{$i as e,Bn as t,Ci as n,Da as r,Hn as i,Kt as a,Pa as o,Qi as s,Sn as c,Vn as l,Wi as u,Zi as d,Zr as f,aa as p,ba as m,br as h,da as g,ea as _,ia as v,la as y,mi as b,na as x,pa as S,ra as C,rr as w,ta as T,vi as E,wa as D,xa as O,xi as k,yi as A,yr as j}from"./router-omaRgp3e.js";import{r as M}from"./index-BXOIXutm.js";var N=E([E(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),A(`spin-container`,`
 position: relative;
 `,[A(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[t()])]),A(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),A(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[k(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),A(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),A(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[k(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=v({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},w.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),i),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:i}=h(e),a=w(`Spin`,`-spin`,N,M,e,t),o=d(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:r},self:i}=a.value,{opacitySpinning:o,color:s,textColor:c}=i;return{"--n-bezier":r,"--n-opacity-spinning":o,"--n-size":typeof t==`number`?b(t):i[n(`size`,t)],"--n-color":s,"--n-text-color":c}}),s=i?j(`spin`,d(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),o,e):void 0,c=f(e,[`spinning`,`show`]),l=r(!1);return m(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{l.value=!0},r),t(()=>{clearTimeout(n)});return}}l.value=c.value}),{mergedClsPrefix:t,active:l,mergedStrokeWidth:d(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return P[typeof n==`number`?`medium`:n]}),cssVars:i?void 0:o,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&p(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?p(`div`,{class:[`${n}-spin-body`,this.themeClass]},p(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):p(`div`,{class:[`${n}-spin-body`,this.themeClass]},p(l,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?p(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},p(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),p(u,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(e,t){return S(),T(`svg`,I,[...t[0]||=[s(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=D({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=v({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let n=t,i=r(Date.now()),l,u=d(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),f=d(()=>{if(u.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-u.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,r=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>r?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),p=d(()=>u.value===null?``:new Date(u.value).toLocaleString(`zh-CN`,{hour12:!1}));return g(()=>{n.compact||(l=setInterval(()=>{i.value=Date.now()},6e4))}),y(()=>{l&&clearInterval(l)}),(n,r)=>{let i=c,l=a;return S(),e(l,null,{trigger:O(()=>[C(i,{size:`small`,type:f.value.type},{default:O(()=>[x(o(f.value.label)+` `,1),!t.compact&&p.value?(S(),T(`span`,z,`· `+o(p.value),1)):_(``,!0)]),_:1},8,[`type`])]),default:O(()=>[s(`div`,null,o(f.value.detail),1),p.value?(S(),T(`div`,B,`数据时间：`+o(p.value),1)):_(``,!0)]),_:1})}}});export{R as n,F as r,V as t};