import{$i as e,$r as t,Bi as n,Cn as r,En as i,Fr as a,Gi as o,Hi as s,Jr as c,Ki as l,Mn as u,Nn as d,Nt as f,On as p,Pn as m,Ta as h,Tn as g,Ui as _,Vi as v,Vn as y,Wi as b,Wn as x,Yi as S,Zi as C,_a as w,ai as T,ar as E,ci as D,cr as O,da as k,ei as A,ir as j,ji as M,jn as N,kn as P,oi as F,on as I,qi as L,qr as R,ra as z,si as B,ta as V,ua as H,ui as U,ur as W,wn as G,zn as K}from"./router-DxH1-4bP.js";import{g as q,h as J,r as Y}from"./index-BYJvH7k2.js";function X(e){let{lineHeight:t,borderRadius:n,fontWeightStrong:r,baseColor:i,dividerColor:a,actionColor:o,textColor1:s,textColor2:l,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,infoColor:h,successColor:g,warningColor:_,errorColor:v,fontSize:y}=e;return Object.assign(Object.assign({},q),{fontSize:y,lineHeight:t,titleFontWeight:r,borderRadius:n,border:`1px solid ${a}`,color:o,titleTextColor:s,iconColor:l,contentTextColor:l,closeBorderRadius:n,closeColorHover:u,closeColorPressed:d,closeIconColor:f,closeIconColorHover:p,closeIconColorPressed:m,borderInfo:`1px solid ${c(i,R(h,{alpha:.25}))}`,colorInfo:c(i,R(h,{alpha:.08})),titleTextColorInfo:s,iconColorInfo:h,contentTextColorInfo:l,closeColorHoverInfo:u,closeColorPressedInfo:d,closeIconColorInfo:f,closeIconColorHoverInfo:p,closeIconColorPressedInfo:m,borderSuccess:`1px solid ${c(i,R(g,{alpha:.25}))}`,colorSuccess:c(i,R(g,{alpha:.08})),titleTextColorSuccess:s,iconColorSuccess:g,contentTextColorSuccess:l,closeColorHoverSuccess:u,closeColorPressedSuccess:d,closeIconColorSuccess:f,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:m,borderWarning:`1px solid ${c(i,R(_,{alpha:.33}))}`,colorWarning:c(i,R(_,{alpha:.08})),titleTextColorWarning:s,iconColorWarning:_,contentTextColorWarning:l,closeColorHoverWarning:u,closeColorPressedWarning:d,closeIconColorWarning:f,closeIconColorHoverWarning:p,closeIconColorPressedWarning:m,borderError:`1px solid ${c(i,R(v,{alpha:.25}))}`,colorError:c(i,R(v,{alpha:.08})),titleTextColorError:s,iconColorError:v,contentTextColorError:l,closeColorHoverError:u,closeColorPressedError:d,closeIconColorError:f,closeIconColorHoverError:p,closeIconColorPressedError:m})}var Z={name:`Alert`,common:r,self:X},Q=F(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[B(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),D(`closable`,[F(`alert-body`,[B(`title`,`
 padding-right: 24px;
 `)])]),B(`icon`,{color:`var(--n-icon-color)`}),F(`alert-body`,{padding:`var(--n-padding)`},[B(`title`,{color:`var(--n-title-text-color)`}),B(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),B(`icon`,`
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
 `),B(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),D(`show-icon`,[F(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),D(`right-adjust`,[F(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),F(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[B(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[T(`& +`,[B(`content`,{marginTop:`9px`})])]),B(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),B(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=L({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},y.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:r,mergedBorderedRef:i,inlineThemeDisabled:a,mergedRtlRef:o}=E(e),s=y(`Alert`,`-alert`,Q,Z,e,r),c=x(`Alert`,o,r),l=n(()=>{let{common:{cubicBezierEaseInOut:n},self:r}=s.value,{fontSize:i,borderRadius:a,titleFontWeight:o,lineHeight:c,iconSize:l,iconMargin:u,iconMarginRtl:d,closeIconSize:f,closeBorderRadius:p,closeSize:m,closeMargin:h,closeMarginRtl:g,padding:_}=r,{type:v}=e,{left:y,right:b}=t(u);return{"--n-bezier":n,"--n-color":r[U(`color`,v)],"--n-close-icon-size":f,"--n-close-border-radius":p,"--n-close-color-hover":r[U(`closeColorHover`,v)],"--n-close-color-pressed":r[U(`closeColorPressed`,v)],"--n-close-icon-color":r[U(`closeIconColor`,v)],"--n-close-icon-color-hover":r[U(`closeIconColorHover`,v)],"--n-close-icon-color-pressed":r[U(`closeIconColorPressed`,v)],"--n-icon-color":r[U(`iconColor`,v)],"--n-border":r[U(`border`,v)],"--n-title-text-color":r[U(`titleTextColor`,v)],"--n-content-text-color":r[U(`contentTextColor`,v)],"--n-line-height":c,"--n-border-radius":a,"--n-font-size":i,"--n-title-font-weight":o,"--n-icon-size":l,"--n-icon-margin":u,"--n-icon-margin-rtl":d,"--n-close-size":m,"--n-close-margin":h,"--n-close-margin-rtl":g,"--n-padding":_,"--n-icon-margin-left":y,"--n-icon-margin-right":b}}),u=a?j(`alert`,n(()=>e.type[0]),l,e):void 0,d=w(!0),f=()=>{let{onAfterLeave:t,onAfterHide:n}=e;t&&t(),n&&n()};return{rtlEnabled:c,mergedClsPrefix:r,mergedBordered:i,visible:d,handleCloseClick:()=>{Promise.resolve(e.onClose?.call(e)).then(e=>{e!==!1&&(d.value=!1)})},handleAfterLeave:()=>{f()},mergedTheme:s,cssVars:a?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),S(p,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:t}=this,n={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?S(`div`,Object.assign({},C(this.$attrs,n)),this.closable&&S(P,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&S(`div`,{class:`${e}-alert__border`}),this.showIcon&&S(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},O(t.icon,()=>[S(K,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return S(u,null);case`info`:return S(d,null);case`warning`:return S(N,null);case`error`:return S(m,null);default:return null}}})])),S(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},W(t.header,t=>{let n=t||this.title;return n?S(`div`,{class:`${e}-alert-body__title`},n):null}),t.default&&S(`div`,{class:`${e}-alert-body__content`},t))):null}})}}),ee=T([T(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),F(`spin-container`,`
 position: relative;
 `,[F(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[G()])]),F(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),F(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[D(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),F(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),F(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[D(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=L({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},y.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),i),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:r}=E(e),i=y(`Spin`,`-spin`,ee,Y,e,t),o=n(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:r}=i.value,{opacitySpinning:a,color:o,textColor:s}=r;return{"--n-bezier":n,"--n-opacity-spinning":a,"--n-size":typeof t==`number`?A(t):r[U(`size`,t)],"--n-color":o,"--n-text-color":s}}),s=r?j(`spin`,n(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),o,e):void 0,c=a(e,[`spinning`,`show`]),l=w(!1);return H(t=>{let n;if(c.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{l.value=!0},r),t(()=>{clearTimeout(n)});return}}l.value=c.value}),{mergedClsPrefix:t,active:l,mergedStrokeWidth:n(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return te[typeof n==`number`?`medium`:n]}),cssVars:r?void 0:o,themeClass:s?.themeClass,onRender:s?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&S(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?S(`div`,{class:[`${n}-spin-body`,this.themeClass]},S(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):S(`div`,{class:[`${n}-spin-body`,this.themeClass]},S(g,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?S(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},S(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),S(M,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=L({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let r=t,i=w(Date.now()),a,c=n(()=>{if(!r.dataAt)return null;let e=new Date(r.dataAt).getTime();return Number.isFinite(e)?e:null}),u=n(()=>{if(c.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-c.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(r.delayedAfterMinutes,1)*60*1e3,n=Math.max(r.expiredAfterMinutes,r.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),d=n(()=>c.value===null?``:new Date(c.value).toLocaleString(`zh-CN`,{hour12:!1}));return V(()=>{r.compact||(a=setInterval(()=>{i.value=Date.now()},6e4))}),e(()=>{a&&clearInterval(a)}),(e,n)=>{let r=I,i=f;return z(),s(i,null,{trigger:k(()=>[l(r,{size:`small`,type:u.value.type},{default:k(()=>[o(h(u.value.label),1),!t.compact&&d.value?(z(),b(`span`,re,`· `+h(d.value),1)):_(``,!0)]),_:1},8,[`type`])]),default:k(()=>[v(`div`,null,h(u.value.detail),1),d.value?(z(),b(`div`,ie,`数据时间：`+h(d.value),1)):_(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};