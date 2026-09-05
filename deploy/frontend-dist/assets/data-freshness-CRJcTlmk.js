import{$n as e,$r as t,Bi as n,Cn as r,Dn as i,Dt as a,En as o,Fi as s,Fn as c,Gi as l,Hr as u,Ii as d,Jr as f,Li as p,Nn as m,On as h,Pi as g,Qn as _,Ri as v,Rn as y,Sn as b,Ti as x,Tn as S,Ui as C,Vi as w,Vr as T,Yi as E,Zi as D,_n as O,aa as k,bn as A,da as j,ei as M,en as N,ia as P,ii as F,ir as I,kr as L,ni as R,nr as z,qi as B,qr as V,ti as H,vn as U,ya as W,yn as G,zi as K}from"./router-CqaEVYAA.js";import{g as q,h as J,r as Y}from"./index-BobVvizI.js";function X(e){let{lineHeight:t,borderRadius:n,fontWeightStrong:r,baseColor:i,dividerColor:a,actionColor:o,textColor1:s,textColor2:c,closeColorHover:l,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,infoColor:h,successColor:g,warningColor:_,errorColor:v,fontSize:y}=e;return Object.assign(Object.assign({},q),{fontSize:y,lineHeight:t,titleFontWeight:r,borderRadius:n,border:`1px solid ${a}`,color:o,titleTextColor:s,iconColor:c,contentTextColor:c,closeBorderRadius:n,closeColorHover:l,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,borderInfo:`1px solid ${u(i,T(h,{alpha:.25}))}`,colorInfo:u(i,T(h,{alpha:.08})),titleTextColorInfo:s,iconColorInfo:h,contentTextColorInfo:c,closeColorHoverInfo:l,closeColorPressedInfo:d,closeIconColorInfo:f,closeIconColorHoverInfo:p,closeIconColorPressedInfo:m,borderSuccess:`1px solid ${u(i,T(g,{alpha:.25}))}`,colorSuccess:u(i,T(g,{alpha:.08})),titleTextColorSuccess:s,iconColorSuccess:g,contentTextColorSuccess:c,closeColorHoverSuccess:l,closeColorPressedSuccess:d,closeIconColorSuccess:f,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:m,borderWarning:`1px solid ${u(i,T(_,{alpha:.33}))}`,colorWarning:u(i,T(_,{alpha:.08})),titleTextColorWarning:s,iconColorWarning:_,contentTextColorWarning:c,closeColorHoverWarning:l,closeColorPressedWarning:d,closeIconColorWarning:f,closeIconColorHoverWarning:p,closeIconColorPressedWarning:m,borderError:`1px solid ${u(i,T(v,{alpha:.25}))}`,colorError:u(i,T(v,{alpha:.08})),titleTextColorError:s,iconColorError:v,contentTextColorError:c,closeColorHoverError:l,closeColorPressedError:d,closeIconColorError:f,closeIconColorHoverError:p,closeIconColorPressedError:m})}var Z={name:`Alert`,common:O,self:X},Q=M(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[H(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),R(`closable`,[M(`alert-body`,[H(`title`,`
 padding-right: 24px;
 `)])]),H(`icon`,{color:`var(--n-icon-color)`}),M(`alert-body`,{padding:`var(--n-padding)`},[H(`title`,{color:`var(--n-title-text-color)`}),H(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),H(`icon`,`
 position: absolute;
 left: 0;
 top: 0;
 align-items: center;
 justify-content: center;
 display: flex;
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 font-size: var(--n-icon-size);
 margin: var(--n-icon-margin);
 `),H(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),R(`show-icon`,[M(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),R(`right-adjust`,[M(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),M(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[H(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[t(`& +`,[H(`content`,{marginTop:`9px`})])]),H(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),H(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=w({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},c.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(t){let{mergedClsPrefixRef:n,mergedBorderedRef:r,inlineThemeDisabled:i,mergedRtlRef:a}=e(t),o=c(`Alert`,`-alert`,Q,Z,t,n),s=y(`Alert`,a,n),l=g(()=>{let{common:{cubicBezierEaseInOut:e},self:n}=o.value,{fontSize:r,borderRadius:i,titleFontWeight:a,lineHeight:s,iconSize:c,iconMargin:l,iconMarginRtl:u,closeIconSize:d,closeBorderRadius:f,closeSize:p,closeMargin:m,closeMarginRtl:h,padding:g}=n,{type:_}=t,{left:v,right:y}=V(l);return{"--n-bezier":e,"--n-color":n[F(`color`,_)],"--n-close-icon-size":d,"--n-close-border-radius":f,"--n-close-color-hover":n[F(`closeColorHover`,_)],"--n-close-color-pressed":n[F(`closeColorPressed`,_)],"--n-close-icon-color":n[F(`closeIconColor`,_)],"--n-close-icon-color-hover":n[F(`closeIconColorHover`,_)],"--n-close-icon-color-pressed":n[F(`closeIconColorPressed`,_)],"--n-icon-color":n[F(`iconColor`,_)],"--n-border":n[F(`border`,_)],"--n-title-text-color":n[F(`titleTextColor`,_)],"--n-content-text-color":n[F(`contentTextColor`,_)],"--n-line-height":s,"--n-border-radius":i,"--n-font-size":r,"--n-title-font-weight":a,"--n-icon-size":c,"--n-icon-margin":l,"--n-icon-margin-rtl":u,"--n-close-size":p,"--n-close-margin":m,"--n-close-margin-rtl":h,"--n-padding":g,"--n-icon-margin-left":v,"--n-icon-margin-right":y}}),u=i?_(`alert`,g(()=>t.type[0]),l,t):void 0,d=j(!0),f=()=>{let{onAfterLeave:e,onAfterHide:n}=t;e&&e(),n&&n()};return{rtlEnabled:s,mergedClsPrefix:n,mergedBordered:r,visible:d,handleCloseClick:()=>{Promise.resolve(t.onClose?.call(t)).then(e=>{e!==!1&&(d.value=!1)})},handleAfterLeave:()=>{f()},mergedTheme:o,cssVars:i?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),C(b,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:t}=this,n={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?C(`div`,Object.assign({},l(this.$attrs,n)),this.closable&&C(r,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&C(`div`,{class:`${e}-alert__border`}),this.showIcon&&C(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},z(t.icon,()=>[C(m,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return C(o,null);case`info`:return C(i,null);case`warning`:return C(S,null);case`error`:return C(h,null);default:return null}}})])),C(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},I(t.header,t=>{let n=t||this.title;return n?C(`div`,{class:`${e}-alert-body__title`},n):null}),t.default&&C(`div`,{class:`${e}-alert-body__content`},t))):null}})}}),ee=t([t(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),M(`spin-container`,`
 position: relative;
 `,[M(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[U()])]),M(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),M(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[R(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),M(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),M(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[R(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=w({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},c.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),A),slots:Object,setup(t){let{mergedClsPrefixRef:n,inlineThemeDisabled:r}=e(t),i=c(`Spin`,`-spin`,ee,Y,t,n),a=g(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:r}=i.value,{opacitySpinning:a,color:o,textColor:s}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof e==`number`?f(e):r[F(`size`,e)],"--n-color":o,"--n-text-color":s}}),o=r?_(`spin`,g(()=>{let{size:e}=t;return typeof e==`number`?String(e):e[0]}),a,t):void 0,s=L(t,[`spinning`,`show`]),l=j(!1);return P(e=>{let n;if(s.value){let{delay:r}=t;if(r){n=window.setTimeout(()=>{l.value=!0},r),e(()=>{clearTimeout(n)});return}}l.value=s.value}),{mergedClsPrefix:n,active:l,mergedStrokeWidth:g(()=>{let{strokeWidth:e}=t;if(e!==void 0)return e;let{size:n}=t;return te[typeof n==`number`?`medium`:n]}),cssVars:r?void 0:a,themeClass:o?.themeClass,onRender:o?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&C(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?C(`div`,{class:[`${n}-spin-body`,this.themeClass]},C(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):C(`div`,{class:[`${n}-spin-body`,this.themeClass]},C(G,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?C(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},C(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),C(x,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=w({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let t=e,r=j(Date.now()),i,o=g(()=>{if(!t.dataAt)return null;let e=new Date(t.dataAt).getTime();return Number.isFinite(e)?e:null}),c=g(()=>{if(o.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=r.value-o.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let n=Math.max(t.delayedAfterMinutes,1)*60*1e3,i=Math.max(t.expiredAfterMinutes,t.delayedAfterMinutes)*60*1e3;return e>i?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(i/6e4)} 分钟未同步`}:e>n?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),l=g(()=>o.value===null?``:new Date(o.value).toLocaleString(`zh-CN`,{hour12:!1}));return E(()=>{t.compact||(i=setInterval(()=>{r.value=Date.now()},6e4))}),B(()=>{i&&clearInterval(i)}),(t,r)=>{let i=N,o=a;return D(),d(o,null,{trigger:k(()=>[n(i,{size:`small`,type:c.value.type},{default:k(()=>[K(W(c.value.label),1),!e.compact&&l.value?(D(),v(`span`,re,`· `+W(l.value),1)):p(``,!0)]),_:1},8,[`type`])]),default:k(()=>[s(`div`,null,W(c.value.detail),1),l.value?(D(),v(`div`,ie,`数据时间：`+W(l.value),1)):p(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};