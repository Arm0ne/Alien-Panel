import{Ai as e,Bi as t,Br as n,Di as r,En as i,Fr as a,Gn as o,Hi as s,Ii as c,Jr as l,Kn as u,Kr as d,Mi as f,Ni as p,Oi as m,On as h,Pi as g,Pr as _,Qi as v,Ri as y,Sr as b,Vr as x,Wi as S,Yn as C,Yr as w,Zi as T,Zn as E,Zr as D,bn as O,dn as k,fn as A,gn as j,hn as M,ji as N,jn as P,ki as F,la as I,na as L,pn as R,qr as z,qt as B,un as V,vi as H,vn as U,xn as W,xt as G,yn as K}from"./router-CBkxnRvy.js";import{g as q,h as J,r as Y}from"./index-BhELkBIb.js";function X(e){let{lineHeight:t,borderRadius:n,fontWeightStrong:r,baseColor:i,dividerColor:o,actionColor:s,textColor1:c,textColor2:l,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,infoColor:h,successColor:g,warningColor:v,errorColor:y,fontSize:b}=e;return Object.assign(Object.assign({},q),{fontSize:b,lineHeight:t,titleFontWeight:r,borderRadius:n,border:`1px solid ${o}`,color:s,titleTextColor:c,iconColor:l,contentTextColor:l,closeBorderRadius:n,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,borderInfo:`1px solid ${a(i,_(h,{alpha:.25}))}`,colorInfo:a(i,_(h,{alpha:.08})),titleTextColorInfo:c,iconColorInfo:h,contentTextColorInfo:l,closeColorHoverInfo:u,closeColorPressedInfo:d,closeIconColorInfo:f,closeIconColorHoverInfo:p,closeIconColorPressedInfo:m,borderSuccess:`1px solid ${a(i,_(g,{alpha:.25}))}`,colorSuccess:a(i,_(g,{alpha:.08})),titleTextColorSuccess:c,iconColorSuccess:g,contentTextColorSuccess:l,closeColorHoverSuccess:u,closeColorPressedSuccess:d,closeIconColorSuccess:f,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:m,borderWarning:`1px solid ${a(i,_(v,{alpha:.33}))}`,colorWarning:a(i,_(v,{alpha:.08})),titleTextColorWarning:c,iconColorWarning:v,contentTextColorWarning:l,closeColorHoverWarning:u,closeColorPressedWarning:d,closeIconColorWarning:f,closeIconColorHoverWarning:p,closeIconColorPressedWarning:m,borderError:`1px solid ${a(i,_(y,{alpha:.25}))}`,colorError:a(i,_(y,{alpha:.08})),titleTextColorError:c,iconColorError:y,contentTextColorError:l,closeColorHoverError:u,closeColorPressedError:d,closeIconColorError:f,closeIconColorHoverError:p,closeIconColorPressedError:m})}var Z={name:`Alert`,common:V,self:X},Q=z(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[l(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),w(`closable`,[z(`alert-body`,[l(`title`,`
 padding-right: 24px;
 `)])]),l(`icon`,{color:`var(--n-icon-color)`}),z(`alert-body`,{padding:`var(--n-padding)`},[l(`title`,{color:`var(--n-title-text-color)`}),l(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),l(`icon`,`
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
 `),l(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),w(`show-icon`,[z(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),w(`right-adjust`,[z(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),z(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[l(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[d(`& +`,[l(`content`,{marginTop:`9px`})])]),l(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),l(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=g({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},h.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:i,inlineThemeDisabled:a,mergedRtlRef:s}=u(e),c=h(`Alert`,`-alert`,Q,Z,e,t),l=P(`Alert`,s,t),d=r(()=>{let{common:{cubicBezierEaseInOut:t},self:r}=c.value,{fontSize:i,borderRadius:a,titleFontWeight:o,lineHeight:s,iconSize:l,iconMargin:u,iconMarginRtl:d,closeIconSize:f,closeBorderRadius:p,closeSize:m,closeMargin:h,closeMarginRtl:g,padding:_}=r,{type:v}=e,{left:y,right:b}=n(u);return{"--n-bezier":t,"--n-color":r[D(`color`,v)],"--n-close-icon-size":f,"--n-close-border-radius":p,"--n-close-color-hover":r[D(`closeColorHover`,v)],"--n-close-color-pressed":r[D(`closeColorPressed`,v)],"--n-close-icon-color":r[D(`closeIconColor`,v)],"--n-close-icon-color-hover":r[D(`closeIconColorHover`,v)],"--n-close-icon-color-pressed":r[D(`closeIconColorPressed`,v)],"--n-icon-color":r[D(`iconColor`,v)],"--n-border":r[D(`border`,v)],"--n-title-text-color":r[D(`titleTextColor`,v)],"--n-content-text-color":r[D(`contentTextColor`,v)],"--n-line-height":s,"--n-border-radius":a,"--n-font-size":i,"--n-title-font-weight":o,"--n-icon-size":l,"--n-icon-margin":u,"--n-icon-margin-rtl":d,"--n-close-size":m,"--n-close-margin":h,"--n-close-margin-rtl":g,"--n-padding":_,"--n-icon-margin-left":y,"--n-icon-margin-right":b}}),f=a?o(`alert`,r(()=>e.type[0]),d,e):void 0,p=L(!0),m=()=>{let{onAfterLeave:t,onAfterHide:n}=e;t&&t(),n&&n()};return{rtlEnabled:l,mergedClsPrefix:t,mergedBordered:i,visible:p,handleCloseClick:()=>{Promise.resolve(e.onClose?.call(e)).then(e=>{e!==!1&&(p.value=!1)})},handleAfterLeave:()=>{m()},mergedTheme:c,cssVars:a?void 0:d,themeClass:f?.themeClass,onRender:f?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),c(M,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:t}=this,n={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?c(`div`,Object.assign({},y(this.$attrs,n)),this.closable&&c(j,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&c(`div`,{class:`${e}-alert__border`}),this.showIcon&&c(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},C(t.icon,()=>[c(i,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return c(K,null);case`info`:return c(O,null);case`warning`:return c(U,null);case`error`:return c(W,null);default:return null}}})])),c(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},E(t.header,t=>{let n=t||this.title;return n?c(`div`,{class:`${e}-alert-body__title`},n):null}),t.default&&c(`div`,{class:`${e}-alert-body__content`},t))):null}})}}),ee=d([d(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),z(`spin-container`,`
 position: relative;
 `,[z(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[k()])]),z(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),z(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[w(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),z(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),z(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[w(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=g({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},h.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),R),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=u(e),i=h(`Spin`,`-spin`,ee,Y,e,t),a=r(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:r}=i.value,{opacitySpinning:a,color:o,textColor:s}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?x(t):r[D(`size`,t)],"--n-color":o,"--n-text-color":s}}),s=n?o(`spin`,r(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),a,e):void 0,c=b(e,[`spinning`,`show`]),l=L(!1);return T(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{l.value=!0},r),t(()=>{clearTimeout(n)});return}}l.value=c.value}),{mergedClsPrefix:t,active:l,mergedStrokeWidth:r(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return te[typeof n==`number`?`medium`:n]}),cssVars:n?void 0:a,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&c(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?c(`div`,{class:[`${n}-spin-body`,this.themeClass]},c(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):c(`div`,{class:[`${n}-spin-body`,this.themeClass]},c(A,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?c(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},c(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),c(H,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=g({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(n){let i=n,a=L(Date.now()),o,c=r(()=>{if(!i.dataAt)return null;let e=new Date(i.dataAt).getTime();return Number.isFinite(e)?e:null}),l=r(()=>{if(c.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=a.value-c.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(i.delayedAfterMinutes,1)*60*1e3,n=Math.max(i.expiredAfterMinutes,i.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),u=r(()=>c.value===null?``:new Date(c.value).toLocaleString(`zh-CN`,{hour12:!1}));return s(()=>{i.compact||(o=setInterval(()=>{a.value=Date.now()},6e4))}),t(()=>{o&&clearInterval(o)}),(t,r)=>{let i=B,a=G;return S(),F(a,null,{trigger:v(()=>[p(i,{size:`small`,type:l.value.type},{default:v(()=>[f(I(l.value.label),1),!n.compact&&u.value?(S(),N(`span`,re,`· `+I(u.value),1)):e(``,!0)]),_:1},8,[`type`])]),default:v(()=>[m(`div`,null,I(l.value.detail),1),u.value?(S(),N(`div`,ie,`数据时间：`+I(u.value),1)):e(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};