import{$i as e,Ai as t,Cr as n,Dn as r,Fi as i,Fr as a,Gi as o,Hr as s,Ir as c,Jr as l,Jt as u,Kn as d,Li as f,Mi as p,Mn as m,Ni as h,Oi as g,Pi as _,Qi as v,Qn as y,Qr as b,Sn as x,St as S,Ui as C,Vi as w,Vr as T,Xn as E,Xr as D,Yr as O,_n as k,bn as A,dn as j,fn as M,gn as N,ji as P,ki as F,kn as I,mn as L,pn as R,qn as z,qr as B,ra as V,ua as H,xn as U,yi as W,yn as G,zi as K}from"./router-Ca5XLGLk.js";import{g as q,h as J,r as Y}from"./index-ChrUCjfM.js";function X(e){let{lineHeight:t,borderRadius:n,fontWeightStrong:r,baseColor:i,dividerColor:o,actionColor:s,textColor1:l,textColor2:u,closeColorHover:d,closeColorPressed:f,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,infoColor:g,successColor:_,warningColor:v,errorColor:y,fontSize:b}=e;return Object.assign(Object.assign({},q),{fontSize:b,lineHeight:t,titleFontWeight:r,borderRadius:n,border:`1px solid ${o}`,color:s,titleTextColor:l,iconColor:u,contentTextColor:u,closeBorderRadius:n,closeColorHover:d,closeColorPressed:f,closeIconColor:p,closeIconColorHover:m,closeIconColorPressed:h,borderInfo:`1px solid ${c(i,a(g,{alpha:.25}))}`,colorInfo:c(i,a(g,{alpha:.08})),titleTextColorInfo:l,iconColorInfo:g,contentTextColorInfo:u,closeColorHoverInfo:d,closeColorPressedInfo:f,closeIconColorInfo:p,closeIconColorHoverInfo:m,closeIconColorPressedInfo:h,borderSuccess:`1px solid ${c(i,a(_,{alpha:.25}))}`,colorSuccess:c(i,a(_,{alpha:.08})),titleTextColorSuccess:l,iconColorSuccess:_,contentTextColorSuccess:u,closeColorHoverSuccess:d,closeColorPressedSuccess:f,closeIconColorSuccess:p,closeIconColorHoverSuccess:m,closeIconColorPressedSuccess:h,borderWarning:`1px solid ${c(i,a(v,{alpha:.33}))}`,colorWarning:c(i,a(v,{alpha:.08})),titleTextColorWarning:l,iconColorWarning:v,contentTextColorWarning:u,closeColorHoverWarning:d,closeColorPressedWarning:f,closeIconColorWarning:p,closeIconColorHoverWarning:m,closeIconColorPressedWarning:h,borderError:`1px solid ${c(i,a(y,{alpha:.25}))}`,colorError:c(i,a(y,{alpha:.08})),titleTextColorError:l,iconColorError:y,contentTextColorError:u,closeColorHoverError:d,closeColorPressedError:f,closeIconColorError:p,closeIconColorHoverError:m,closeIconColorPressedError:h})}var Z={name:`Alert`,common:j,self:X},Q=l(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[O(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),D(`closable`,[l(`alert-body`,[O(`title`,`
 padding-right: 24px;
 `)])]),O(`icon`,{color:`var(--n-icon-color)`}),l(`alert-body`,{padding:`var(--n-padding)`},[O(`title`,{color:`var(--n-title-text-color)`}),O(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),O(`icon`,`
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
 `),O(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),D(`show-icon`,[l(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),D(`right-adjust`,[l(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),l(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[O(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[B(`& +`,[O(`content`,{marginTop:`9px`})])]),O(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),O(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=i({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},I.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,inlineThemeDisabled:r,mergedRtlRef:i}=z(e),a=I(`Alert`,`-alert`,Q,Z,e,t),o=m(`Alert`,i,t),s=g(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=a.value,{fontSize:r,borderRadius:i,titleFontWeight:o,lineHeight:s,iconSize:c,iconMargin:l,iconMarginRtl:u,closeIconSize:d,closeBorderRadius:f,closeSize:p,closeMargin:m,closeMarginRtl:h,padding:g}=n,{type:_}=e,{left:v,right:y}=T(l);return{"--n-bezier":t,"--n-color":n[b(`color`,_)],"--n-close-icon-size":d,"--n-close-border-radius":f,"--n-close-color-hover":n[b(`closeColorHover`,_)],"--n-close-color-pressed":n[b(`closeColorPressed`,_)],"--n-close-icon-color":n[b(`closeIconColor`,_)],"--n-close-icon-color-hover":n[b(`closeIconColorHover`,_)],"--n-close-icon-color-pressed":n[b(`closeIconColorPressed`,_)],"--n-icon-color":n[b(`iconColor`,_)],"--n-border":n[b(`border`,_)],"--n-title-text-color":n[b(`titleTextColor`,_)],"--n-content-text-color":n[b(`contentTextColor`,_)],"--n-line-height":s,"--n-border-radius":i,"--n-font-size":r,"--n-title-font-weight":o,"--n-icon-size":c,"--n-icon-margin":l,"--n-icon-margin-rtl":u,"--n-close-size":p,"--n-close-margin":m,"--n-close-margin-rtl":h,"--n-padding":g,"--n-icon-margin-left":v,"--n-icon-margin-right":y}}),c=r?d(`alert`,g(()=>e.type[0]),s,e):void 0,l=V(!0),u=()=>{let{onAfterLeave:t,onAfterHide:n}=e;t&&t(),n&&n()};return{rtlEnabled:o,mergedClsPrefix:t,mergedBordered:n,visible:l,handleCloseClick:()=>{Promise.resolve(e.onClose?.call(e)).then(e=>{e!==!1&&(l.value=!1)})},handleAfterLeave:()=>{u()},mergedTheme:a,cssVars:r?void 0:s,themeClass:c?.themeClass,onRender:c?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),f(N,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:t}=this,n={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?f(`div`,Object.assign({},K(this.$attrs,n)),this.closable&&f(k,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&f(`div`,{class:`${e}-alert__border`}),this.showIcon&&f(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},E(t.icon,()=>[f(r,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return f(A,null);case`info`:return f(U,null);case`warning`:return f(G,null);case`error`:return f(x,null);default:return null}}})])),f(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},y(t.header,t=>{let n=t||this.title;return n?f(`div`,{class:`${e}-alert-body__title`},n):null}),t.default&&f(`div`,{class:`${e}-alert-body__content`},t))):null}})}}),ee=B([B(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),l(`spin-container`,`
 position: relative;
 `,[l(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[M()])]),l(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),l(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[D(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),l(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),l(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[D(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=i({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},I.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),L),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:r}=z(e),i=I(`Spin`,`-spin`,ee,Y,e,t),a=g(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:r}=i.value,{opacitySpinning:a,color:o,textColor:c}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?s(t):r[b(`size`,t)],"--n-color":o,"--n-text-color":c}}),o=r?d(`spin`,g(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),a,e):void 0,c=n(e,[`spinning`,`show`]),l=V(!1);return v(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{l.value=!0},r),t(()=>{clearTimeout(n)});return}}l.value=c.value}),{mergedClsPrefix:t,active:l,mergedStrokeWidth:g(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return te[typeof n==`number`?`medium`:n]}),cssVars:r?void 0:a,themeClass:o?.themeClass,onRender:o?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&f(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?f(`div`,{class:[`${n}-spin-body`,this.themeClass]},f(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):f(`div`,{class:[`${n}-spin-body`,this.themeClass]},f(R,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?f(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},f(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),f(W,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=i({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(n){let r=n,i=V(Date.now()),a,s=g(()=>{if(!r.dataAt)return null;let e=new Date(r.dataAt).getTime();return Number.isFinite(e)?e:null}),c=g(()=>{if(s.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-s.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(r.delayedAfterMinutes,1)*60*1e3,n=Math.max(r.expiredAfterMinutes,r.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),l=g(()=>s.value===null?``:new Date(s.value).toLocaleString(`zh-CN`,{hour12:!1}));return C(()=>{r.compact||(a=setInterval(()=>{i.value=Date.now()},6e4))}),w(()=>{a&&clearInterval(a)}),(r,i)=>{let a=u,s=S;return o(),t(s,null,{trigger:e(()=>[_(a,{size:`small`,type:c.value.type},{default:e(()=>[h(H(c.value.label),1),!n.compact&&l.value?(o(),p(`span`,re,`· `+H(l.value),1)):P(``,!0)]),_:1},8,[`type`])]),default:e(()=>[F(`div`,null,H(c.value.detail),1),l.value?(o(),p(`div`,ie,`数据时间：`+H(l.value),1)):P(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};