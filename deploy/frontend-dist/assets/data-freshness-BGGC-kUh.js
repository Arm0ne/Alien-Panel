import{$i as e,$r as t,Ci as n,Fa as r,Hn as i,Jt as a,Oa as o,Qi as s,Sa as c,Sr as l,Ta as u,Ti as d,Un as f,Wi as p,Wn as m,Zi as h,aa as g,ar as _,bi as v,ea as y,fa as b,gi as x,ia as S,ma as C,oa as w,ra as T,ta as E,ua as D,wn as O,xa as k,xi as A,xr as j}from"./router-XnvCoRTR.js";import{r as M}from"./index-CjMER2w_.js";var N=v([v(`@keyframes spin-rotate`,`
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
 `,[i()])]),A(`spin-body`,`
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
 `,[n(`rotate`,`
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
 `,[n(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=g({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},_.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),m),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:r}=l(e),i=_(`Spin`,`-spin`,N,M,e,n),a=h(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:r}=i.value,{opacitySpinning:a,color:o,textColor:s}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?x(t):r[d(`size`,t)],"--n-color":o,"--n-text-color":s}}),s=r?j(`spin`,h(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),a,e):void 0,c=t(e,[`spinning`,`show`]),u=o(!1);return k(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{u.value=!0},r),t(()=>{clearTimeout(n)});return}}u.value=c.value}),{mergedClsPrefix:n,active:u,mergedStrokeWidth:h(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return P[typeof n==`number`?`medium`:n]}),cssVars:r?void 0:a,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&w(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?w(`div`,{class:[`${n}-spin-body`,this.themeClass]},w(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):w(`div`,{class:[`${n}-spin-body`,this.themeClass]},w(f,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?w(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},w(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),w(p,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(e,t){return C(),E(`svg`,I,[...t[0]||=[s(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=u({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=g({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let n=t,i=o(Date.now()),l,u=h(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),d=h(()=>{if(u.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-u.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,r=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>r?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),f=h(()=>u.value===null?``:new Date(u.value).toLocaleString(`zh-CN`,{hour12:!1}));return b(()=>{n.compact||(l=setInterval(()=>{i.value=Date.now()},6e4))}),D(()=>{l&&clearInterval(l)}),(n,i)=>{let o=O,l=a;return C(),e(l,null,{trigger:c(()=>[S(o,{size:`small`,type:d.value.type},{default:c(()=>[T(r(d.value.label)+` `,1),!t.compact&&f.value?(C(),E(`span`,z,`· `+r(f.value),1)):y(``,!0)]),_:1},8,[`type`])]),default:c(()=>[s(`div`,null,r(d.value.detail),1),f.value?(C(),E(`div`,B,`数据时间：`+r(f.value),1)):y(``,!0)]),_:1})}}});export{R as n,F as r,V as t};