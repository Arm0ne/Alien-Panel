import{$r as e,Br as t,Ci as n,Cn as r,Dr as i,En as a,Fi as o,Gi as s,Gr as c,Ii as l,In as u,Kr as d,Li as f,Mi as p,Ni as m,Nn as h,Pi as g,Qr as _,Qt as v,Ri as y,Tn as b,Tt as x,Ui as S,Vi as C,Xn as w,Yi as T,Zn as E,Zr as D,_a as O,_n as k,bn as A,ei as j,er as M,gn as N,hn as P,jn as F,la as I,na as L,ni as R,nr as z,qi as B,ra as V,vn as H,wn as U,xn as W,zi as G,zr as K}from"./router-BXCsM_PA.js";import{g as q,h as J,r as Y}from"./index-BCfmTlJD.js";function X(e){let{lineHeight:n,borderRadius:r,fontWeightStrong:i,baseColor:a,dividerColor:o,actionColor:s,textColor1:c,textColor2:l,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,infoColor:h,successColor:g,warningColor:_,errorColor:v,fontSize:y}=e;return Object.assign(Object.assign({},q),{fontSize:y,lineHeight:n,titleFontWeight:i,borderRadius:r,border:`1px solid ${o}`,color:s,titleTextColor:c,iconColor:l,contentTextColor:l,closeBorderRadius:r,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,borderInfo:`1px solid ${t(a,K(h,{alpha:.25}))}`,colorInfo:t(a,K(h,{alpha:.08})),titleTextColorInfo:c,iconColorInfo:h,contentTextColorInfo:l,closeColorHoverInfo:u,closeColorPressedInfo:d,closeIconColorInfo:f,closeIconColorHoverInfo:p,closeIconColorPressedInfo:m,borderSuccess:`1px solid ${t(a,K(g,{alpha:.25}))}`,colorSuccess:t(a,K(g,{alpha:.08})),titleTextColorSuccess:c,iconColorSuccess:g,contentTextColorSuccess:l,closeColorHoverSuccess:u,closeColorPressedSuccess:d,closeIconColorSuccess:f,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:m,borderWarning:`1px solid ${t(a,K(_,{alpha:.33}))}`,colorWarning:t(a,K(_,{alpha:.08})),titleTextColorWarning:c,iconColorWarning:_,contentTextColorWarning:l,closeColorHoverWarning:u,closeColorPressedWarning:d,closeIconColorWarning:f,closeIconColorHoverWarning:p,closeIconColorPressedWarning:m,borderError:`1px solid ${t(a,K(v,{alpha:.25}))}`,colorError:t(a,K(v,{alpha:.08})),titleTextColorError:c,iconColorError:v,contentTextColorError:l,closeColorHoverError:u,closeColorPressedError:d,closeIconColorError:f,closeIconColorHoverError:p,closeIconColorPressedError:m})}var Z={name:`Alert`,common:P,self:X},Q=_(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[e(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),j(`closable`,[_(`alert-body`,[e(`title`,`
 padding-right: 24px;
 `)])]),e(`icon`,{color:`var(--n-icon-color)`}),_(`alert-body`,{padding:`var(--n-padding)`},[e(`title`,{color:`var(--n-title-text-color)`}),e(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),e(`icon`,`
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
 `),e(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),j(`show-icon`,[_(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),j(`right-adjust`,[_(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),_(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[e(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[D(`& +`,[e(`content`,{marginTop:`9px`})])]),e(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),e(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=G({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},h.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,inlineThemeDisabled:r,mergedRtlRef:i}=E(e),a=h(`Alert`,`-alert`,Q,Z,e,t),o=u(`Alert`,i,t),s=p(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=a.value,{fontSize:r,borderRadius:i,titleFontWeight:o,lineHeight:s,iconSize:l,iconMargin:u,iconMarginRtl:d,closeIconSize:f,closeBorderRadius:p,closeSize:m,closeMargin:h,closeMarginRtl:g,padding:_}=n,{type:v}=e,{left:y,right:b}=c(u);return{"--n-bezier":t,"--n-color":n[R(`color`,v)],"--n-close-icon-size":f,"--n-close-border-radius":p,"--n-close-color-hover":n[R(`closeColorHover`,v)],"--n-close-color-pressed":n[R(`closeColorPressed`,v)],"--n-close-icon-color":n[R(`closeIconColor`,v)],"--n-close-icon-color-hover":n[R(`closeIconColorHover`,v)],"--n-close-icon-color-pressed":n[R(`closeIconColorPressed`,v)],"--n-icon-color":n[R(`iconColor`,v)],"--n-border":n[R(`border`,v)],"--n-title-text-color":n[R(`titleTextColor`,v)],"--n-content-text-color":n[R(`contentTextColor`,v)],"--n-line-height":s,"--n-border-radius":i,"--n-font-size":r,"--n-title-font-weight":o,"--n-icon-size":l,"--n-icon-margin":u,"--n-icon-margin-rtl":d,"--n-close-size":m,"--n-close-margin":h,"--n-close-margin-rtl":g,"--n-padding":_,"--n-icon-margin-left":y,"--n-icon-margin-right":b}}),l=r?w(`alert`,p(()=>e.type[0]),s,e):void 0,d=I(!0),f=()=>{let{onAfterLeave:t,onAfterHide:n}=e;t&&t(),n&&n()};return{rtlEnabled:o,mergedClsPrefix:t,mergedBordered:n,visible:d,handleCloseClick:()=>{Promise.resolve(e.onClose?.call(e)).then(e=>{e!==!1&&(d.value=!1)})},handleAfterLeave:()=>{f()},mergedTheme:a,cssVars:r?void 0:s,themeClass:l?.themeClass,onRender:l?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),C(A,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:t}=this,n={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?C(`div`,Object.assign({},S(this.$attrs,n)),this.closable&&C(W,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&C(`div`,{class:`${e}-alert__border`}),this.showIcon&&C(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},M(t.icon,()=>[C(F,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return C(U,null);case`info`:return C(b,null);case`warning`:return C(r,null);case`error`:return C(a,null);default:return null}}})])),C(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},z(t.header,t=>{let n=t||this.title;return n?C(`div`,{class:`${e}-alert-body__title`},n):null}),t.default&&C(`div`,{class:`${e}-alert-body__content`},t))):null}})}}),ee=D([D(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),_(`spin-container`,`
 position: relative;
 `,[_(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[N()])]),_(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),_(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[j(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),_(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),_(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[j(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=G({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},h.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),H),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=E(e),r=h(`Spin`,`-spin`,ee,Y,e,t),a=p(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{opacitySpinning:a,color:o,textColor:s}=i;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?d(t):i[R(`size`,t)],"--n-color":o,"--n-text-color":s}}),o=n?w(`spin`,p(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),a,e):void 0,s=i(e,[`spinning`,`show`]),c=I(!1);return L(t=>{let n;if(s.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{c.value=!0},r),t(()=>{clearTimeout(n)});return}}c.value=s.value}),{mergedClsPrefix:t,active:c,mergedStrokeWidth:p(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return te[typeof n==`number`?`medium`:n]}),cssVars:n?void 0:a,themeClass:o?.themeClass,onRender:o?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:r,description:i}=this,a=t.icon&&this.rotate,o=(i||t.description)&&C(`div`,{class:`${r}-spin-description`},i||t.description?.call(t)),s=t.icon?C(`div`,{class:[`${r}-spin-body`,this.themeClass]},C(`div`,{class:[`${r}-spin`,a&&`${r}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),o):C(`div`,{class:[`${r}-spin-body`,this.themeClass]},C(k,{clsPrefix:r,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${r}-spin`}),o);return(e=this.onRender)==null||e.call(this),t.default?C(`div`,{class:[`${r}-spin-container`,this.themeClass],style:this.cssVars},C(`div`,{class:[`${r}-spin-content`,this.active&&`${r}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),C(n,{name:`fade-in-transition`},{default:()=>this.active?s:null})):s}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=G({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let t=e,n=I(Date.now()),r,i=p(()=>{if(!t.dataAt)return null;let e=new Date(t.dataAt).getTime();return Number.isFinite(e)?e:null}),a=p(()=>{if(i.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=n.value-i.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let r=Math.max(t.delayedAfterMinutes,1)*60*1e3,a=Math.max(t.expiredAfterMinutes,t.delayedAfterMinutes)*60*1e3;return e>a?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(a/6e4)} 分钟未同步`}:e>r?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),c=p(()=>i.value===null?``:new Date(i.value).toLocaleString(`zh-CN`,{hour12:!1}));return B(()=>{t.compact||(r=setInterval(()=>{n.value=Date.now()},6e4))}),s(()=>{r&&clearInterval(r)}),(t,n)=>{let r=v,i=x;return T(),g(i,null,{trigger:V(()=>[y(r,{size:`small`,type:a.value.type},{default:V(()=>[f(O(a.value.label),1),!e.compact&&c.value?(T(),l(`span`,re,`· `+O(c.value),1)):o(``,!0)]),_:1},8,[`type`])]),default:V(()=>[m(`div`,null,O(a.value.detail),1),c.value?(T(),l(`div`,ie,`数据时间：`+O(c.value),1)):o(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};