import{Aa as e,Ca as t,Cr as n,Da as r,Ei as i,Gn as a,La as o,Si as s,Sr as c,Tn as l,Un as u,Wn as d,Yt as f,_i as p,aa as m,ca as h,ea as g,ei as _,fa as v,ga as y,ia as b,ma as x,na as S,oa as C,or as w,qi as T,ra as E,sa as D,ta as O,wa as k,wi as A,xi as j}from"./router-iRgl5rSe.js";import{r as M}from"./index--kMbNr0h.js";var N=j([j(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),s(`spin-container`,`
 position: relative;
 `,[s(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[u()])]),s(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),s(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[A(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),s(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),s(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[A(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),P={small:20,medium:18,large:16},F=D({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},w.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),a),slots:Object,setup(r){let{mergedClsPrefixRef:a,inlineThemeDisabled:o}=n(r),s=w(`Spin`,`-spin`,N,M,r,a),l=g(()=>{let{size:e}=r,{common:{cubicBezierEaseInOut:t},self:n}=s.value,{opacitySpinning:a,color:o,textColor:c}=n;return{"--n-bezier":t,"--n-opacity-spinning":a,"--n-size":typeof e==`number`?p(e):n[i(`size`,e)],"--n-color":o,"--n-text-color":c}}),u=o?c(`spin`,g(()=>{let{size:e}=r;return typeof e==`number`?String(e):e[0]}),l,r):void 0,d=_(r,[`spinning`,`show`]),f=e(!1);return t(e=>{let t;if(d.value){let{delay:n}=r;if(n){t=window.setTimeout(()=>{f.value=!0},n),e(()=>{clearTimeout(t)});return}}f.value=d.value}),{mergedClsPrefix:a,active:f,mergedStrokeWidth:g(()=>{let{strokeWidth:e}=r;if(e!==void 0)return e;let{size:t}=r;return P[typeof t==`number`?`medium`:t]}),cssVars:o?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&h(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?h(`div`,{class:[`${n}-spin-body`,this.themeClass]},h(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):h(`div`,{class:[`${n}-spin-body`,this.themeClass]},h(d,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?h(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},h(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),h(T,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),I={class:`inline-block`,viewBox:`0 0 24 24`,width:`1em`,height:`1em`};function L(e,t){return y(),b(`svg`,I,[...t[0]||=[O(`path`,{fill:`currentColor`,d:`M17.65 6.35A7.96 7.96 0 0 0 12 4a8 8 0 0 0-8 8a8 8 0 0 0 8 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18a6 6 0 0 1-6-6a6 6 0 0 1 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4z`},null,-1)]])}var R=r({name:`mdi-refresh`,render:L}),z={key:0,class:`ml-4px`},B={key:0,class:`mt-2px text-12px opacity-80`},V=D({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let n=t,r=e(Date.now()),i,a=g(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),s=g(()=>{if(a.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=r.value-a.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,i=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>i?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(i/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),c=g(()=>a.value===null?``:new Date(a.value).toLocaleString(`zh-CN`,{hour12:!1}));return x(()=>{n.compact||(i=setInterval(()=>{r.value=Date.now()},6e4))}),v(()=>{i&&clearInterval(i)}),(e,n)=>{let r=l,i=f;return y(),S(i,null,{trigger:k(()=>[C(r,{size:`small`,type:s.value.type},{default:k(()=>[m(o(s.value.label)+` `,1),!t.compact&&c.value?(y(),b(`span`,z,`· `+o(c.value),1)):E(``,!0)]),_:1},8,[`type`])]),default:k(()=>[O(`div`,null,o(s.value.detail),1),c.value?(y(),b(`div`,B,`数据时间：`+o(c.value),1)):E(``,!0)]),_:1})}}});export{R as n,F as r,V as t};