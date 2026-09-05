import{$i as e,$n as t,$r as n,Ai as r,An as i,Bi as a,Cn as o,Ct as s,Fi as c,Hi as l,Hr as u,Ii as d,Ir as f,Jn as p,Jr as m,Ki as h,Lr as g,Mi as _,Ni as v,Nn as y,On as b,Pi as x,Ri as S,Sn as C,Ur as w,Wi as T,Xr as E,Yr as D,Yt as O,Zn as k,Zr as A,_n as j,bi as M,bn as N,da as P,ea as F,fn as I,hn as L,ia as R,ji as z,ki as B,mn as V,pn as H,qn as U,vn as W,wr as G,xn as K}from"./router-CLsq84Q3.js";import{g as q,h as J,r as Y}from"./index-ChGrRyuT.js";function X(e){let{lineHeight:t,borderRadius:n,fontWeightStrong:r,baseColor:i,dividerColor:a,actionColor:o,textColor1:s,textColor2:c,closeColorHover:l,closeColorPressed:u,closeIconColor:d,closeIconColorHover:p,closeIconColorPressed:m,infoColor:h,successColor:_,warningColor:v,errorColor:y,fontSize:b}=e;return Object.assign(Object.assign({},q),{fontSize:b,lineHeight:t,titleFontWeight:r,borderRadius:n,border:`1px solid ${a}`,color:o,titleTextColor:s,iconColor:c,contentTextColor:c,closeBorderRadius:n,closeColorHover:l,closeColorPressed:u,closeIconColor:d,closeIconColorHover:p,closeIconColorPressed:m,borderInfo:`1px solid ${g(i,f(h,{alpha:.25}))}`,colorInfo:g(i,f(h,{alpha:.08})),titleTextColorInfo:s,iconColorInfo:h,contentTextColorInfo:c,closeColorHoverInfo:l,closeColorPressedInfo:u,closeIconColorInfo:d,closeIconColorHoverInfo:p,closeIconColorPressedInfo:m,borderSuccess:`1px solid ${g(i,f(_,{alpha:.25}))}`,colorSuccess:g(i,f(_,{alpha:.08})),titleTextColorSuccess:s,iconColorSuccess:_,contentTextColorSuccess:c,closeColorHoverSuccess:l,closeColorPressedSuccess:u,closeIconColorSuccess:d,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:m,borderWarning:`1px solid ${g(i,f(v,{alpha:.33}))}`,colorWarning:g(i,f(v,{alpha:.08})),titleTextColorWarning:s,iconColorWarning:v,contentTextColorWarning:c,closeColorHoverWarning:l,closeColorPressedWarning:u,closeIconColorWarning:d,closeIconColorHoverWarning:p,closeIconColorPressedWarning:m,borderError:`1px solid ${g(i,f(y,{alpha:.25}))}`,colorError:g(i,f(y,{alpha:.08})),titleTextColorError:s,iconColorError:y,contentTextColorError:c,closeColorHoverError:l,closeColorPressedError:u,closeIconColorError:d,closeIconColorHoverError:p,closeIconColorPressedError:m})}var Z={name:`Alert`,common:I,self:X},Q=D(`alert`,`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[E(`border`,`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),A(`closable`,[D(`alert-body`,[E(`title`,`
 padding-right: 24px;
 `)])]),E(`icon`,{color:`var(--n-icon-color)`}),D(`alert-body`,{padding:`var(--n-padding)`},[E(`title`,{color:`var(--n-title-text-color)`}),E(`content`,{color:`var(--n-content-text-color)`})]),J({originalTransition:`transform .3s var(--n-bezier)`,enterToProps:{transform:`scale(1)`},leaveToProps:{transform:`scale(0.9)`}}),E(`icon`,`
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
 `),E(`close`,`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),A(`show-icon`,[D(`alert-body`,{paddingLeft:`calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))`})]),A(`right-adjust`,[D(`alert-body`,{paddingRight:`calc(var(--n-close-size) + var(--n-padding) + 2px)`})]),D(`alert-body`,`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[E(`title`,`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[m(`& +`,[E(`content`,{marginTop:`9px`})])]),E(`content`,{transition:`color .3s var(--n-bezier)`,fontSize:`var(--n-font-size)`})]),E(`icon`,{transition:`color .3s var(--n-bezier)`})]),$=d({name:`Alert`,inheritAttrs:!1,props:Object.assign(Object.assign({},i.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:`default`},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:r,inlineThemeDisabled:a,mergedRtlRef:o}=p(e),s=i(`Alert`,`-alert`,Q,Z,e,t),c=y(`Alert`,o,t),l=B(()=>{let{common:{cubicBezierEaseInOut:t},self:r}=s.value,{fontSize:i,borderRadius:a,titleFontWeight:o,lineHeight:c,iconSize:l,iconMargin:d,iconMarginRtl:f,closeIconSize:p,closeBorderRadius:m,closeSize:h,closeMargin:g,closeMarginRtl:_,padding:v}=r,{type:y}=e,{left:b,right:x}=u(d);return{"--n-bezier":t,"--n-color":r[n(`color`,y)],"--n-close-icon-size":p,"--n-close-border-radius":m,"--n-close-color-hover":r[n(`closeColorHover`,y)],"--n-close-color-pressed":r[n(`closeColorPressed`,y)],"--n-close-icon-color":r[n(`closeIconColor`,y)],"--n-close-icon-color-hover":r[n(`closeIconColorHover`,y)],"--n-close-icon-color-pressed":r[n(`closeIconColorPressed`,y)],"--n-icon-color":r[n(`iconColor`,y)],"--n-border":r[n(`border`,y)],"--n-title-text-color":r[n(`titleTextColor`,y)],"--n-content-text-color":r[n(`contentTextColor`,y)],"--n-line-height":c,"--n-border-radius":a,"--n-font-size":i,"--n-title-font-weight":o,"--n-icon-size":l,"--n-icon-margin":d,"--n-icon-margin-rtl":f,"--n-close-size":h,"--n-close-margin":g,"--n-close-margin-rtl":_,"--n-padding":v,"--n-icon-margin-left":b,"--n-icon-margin-right":x}}),d=a?U(`alert`,B(()=>e.type[0]),l,e):void 0,f=R(!0),m=()=>{let{onAfterLeave:t,onAfterHide:n}=e;t&&t(),n&&n()};return{rtlEnabled:c,mergedClsPrefix:t,mergedBordered:r,visible:f,handleCloseClick:()=>{Promise.resolve(e.onClose?.call(e)).then(e=>{e!==!1&&(f.value=!1)})},handleAfterLeave:()=>{m()},mergedTheme:s,cssVars:a?void 0:l,themeClass:d?.themeClass,onRender:d?.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),S(j,{onAfterLeave:this.handleAfterLeave},{default:()=>{let{mergedClsPrefix:e,$slots:n}=this,r={class:[`${e}-alert`,this.themeClass,this.closable&&`${e}-alert--closable`,this.showIcon&&`${e}-alert--show-icon`,!this.title&&this.closable&&`${e}-alert--right-adjust`,this.rtlEnabled&&`${e}-alert--rtl`],style:this.cssVars,role:`alert`};return this.visible?S(`div`,Object.assign({},a(this.$attrs,r)),this.closable&&S(W,{clsPrefix:e,class:`${e}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&S(`div`,{class:`${e}-alert__border`}),this.showIcon&&S(`div`,{class:`${e}-alert__icon`,"aria-hidden":`true`},k(n.icon,()=>[S(b,{clsPrefix:e},{default:()=>{switch(this.type){case`success`:return S(K,null);case`info`:return S(C,null);case`warning`:return S(N,null);case`error`:return S(o,null);default:return null}}})])),S(`div`,{class:[`${e}-alert-body`,this.mergedBordered&&`${e}-alert-body--bordered`]},t(n.header,t=>{let n=t||this.title;return n?S(`div`,{class:`${e}-alert-body__title`},n):null}),n.default&&S(`div`,{class:`${e}-alert-body__content`},n))):null}})}}),ee=m([m(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),D(`spin-container`,`
 position: relative;
 `,[D(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[H()])]),D(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),D(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[A(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),D(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),D(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[A(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),te={small:20,medium:18,large:16},ne=d({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},i.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),L),slots:Object,setup(t){let{mergedClsPrefixRef:r,inlineThemeDisabled:a}=p(t),o=i(`Spin`,`-spin`,ee,Y,t,r),s=B(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:r},self:i}=o.value,{opacitySpinning:a,color:s,textColor:c}=i;return{"--n-bezier":r,"--n-opacity-spinning":a,"--n-size":typeof e==`number`?w(e):i[n(`size`,e)],"--n-color":s,"--n-text-color":c}}),c=a?U(`spin`,B(()=>{let{size:e}=t;return typeof e==`number`?String(e):e[0]}),s,t):void 0,l=G(t,[`spinning`,`show`]),u=R(!1);return e(e=>{let n;if(l.value){let{delay:r}=t;if(r){n=window.setTimeout(()=>{u.value=!0},r),e(()=>{clearTimeout(n)});return}}u.value=l.value}),{mergedClsPrefix:r,active:u,mergedStrokeWidth:B(()=>{let{strokeWidth:e}=t;if(e!==void 0)return e;let{size:n}=t;return te[typeof n==`number`?`medium`:n]}),cssVars:a?void 0:s,themeClass:c?.themeClass,onRender:c?.onRender}},render(){var e;let{$slots:t,mergedClsPrefix:n,description:r}=this,i=t.icon&&this.rotate,a=(r||t.description)&&S(`div`,{class:`${n}-spin-description`},r||t.description?.call(t)),o=t.icon?S(`div`,{class:[`${n}-spin-body`,this.themeClass]},S(`div`,{class:[`${n}-spin`,i&&`${n}-spin--rotate`],style:t.default?``:this.cssVars},t.icon()),a):S(`div`,{class:[`${n}-spin-body`,this.themeClass]},S(V,{clsPrefix:n,style:t.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${n}-spin`}),a);return(e=this.onRender)==null||e.call(this),t.default?S(`div`,{class:[`${n}-spin-container`,this.themeClass],style:this.cssVars},S(`div`,{class:[`${n}-spin-content`,this.active&&`${n}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),S(M,{name:`fade-in-transition`},{default:()=>this.active?o:null})):o}}),re={key:0,class:`ml-4px`},ie={key:0,class:`mt-2px text-12px opacity-80`},ae=d({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let t=e,n=R(Date.now()),i,a=B(()=>{if(!t.dataAt)return null;let e=new Date(t.dataAt).getTime();return Number.isFinite(e)?e:null}),o=B(()=>{if(a.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=n.value-a.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let r=Math.max(t.delayedAfterMinutes,1)*60*1e3,i=Math.max(t.expiredAfterMinutes,t.delayedAfterMinutes)*60*1e3;return e>i?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(i/6e4)} 分钟未同步`}:e>r?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),u=B(()=>a.value===null?``:new Date(a.value).toLocaleString(`zh-CN`,{hour12:!1}));return T(()=>{t.compact||(i=setInterval(()=>{n.value=Date.now()},6e4))}),l(()=>{i&&clearInterval(i)}),(t,n)=>{let i=O,a=s;return h(),z(a,null,{trigger:F(()=>[c(i,{size:`small`,type:o.value.type},{default:F(()=>[x(P(o.value.label),1),!e.compact&&u.value?(h(),v(`span`,re,`· `+P(u.value),1)):_(``,!0)]),_:1},8,[`type`])]),default:F(()=>[r(`div`,null,P(o.value.detail),1),u.value?(h(),v(`div`,ie,`数据时间：`+P(u.value),1)):_(``,!0)]),_:1})}}});export{ne as n,$ as r,ae as t};