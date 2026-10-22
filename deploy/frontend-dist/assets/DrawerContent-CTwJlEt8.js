import{$i as e,Bi as t,Bn as n,Cn as r,Dr as i,En as a,Er as o,Fr as s,Ji as c,Li as l,Or as u,Pr as d,Ri as f,Sr as p,Tn as m,Tr as h,Ui as g,_r as _,ar as v,dn as y,ea as b,ei as x,er as S,ki as C,la as w,lr as T,na as E,ni as D,pr as O,ri as k,rn as A,sa as j,sn as M,sr as N,ti as P,tr as F,vr as I,wi as L,xi as R,xr as z,zi as B,zn as V}from"./router-BEJlwQm6.js";import{r as H}from"./index-DK8sHnu3.js";var U=l({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(t){let r=j(!!t.show),s=j(null),l=B(u),d=0,f=``,p=null,m=j(!1),_=j(!1),v=C(()=>t.placement===`top`||t.placement===`bottom`),{mergedClsPrefixRef:y,mergedRtlRef:x}=n(t),S=a(`Drawer`,x,y),w=P,T=e=>{_.value=!0,d=v.value?e.clientY:e.clientX,f=document.body.style.cursor,document.body.style.cursor=v.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,N),document.body.addEventListener(`mouseleave`,w),document.body.addEventListener(`mouseup`,P)},E=()=>{p!==null&&(window.clearTimeout(p),p=null),_.value?m.value=!0:p=window.setTimeout(()=>{m.value=!0},300)},D=()=>{p!==null&&(window.clearTimeout(p),p=null),m.value=!1},{doUpdateHeight:O,doUpdateWidth:k}=l,A=e=>{let{maxWidth:n}=t;if(n&&e>n)return n;let{minWidth:r}=t;return r&&e<r?r:e},M=e=>{let{maxHeight:n}=t;if(n&&e>n)return n;let{minHeight:r}=t;return r&&e<r?r:e};function N(e){if(_.value)if(v.value){let n=s.value?.offsetHeight||0,r=d-e.clientY;n+=t.placement===`bottom`?r:-r,n=M(n),O(n),d=e.clientY}else{let n=s.value?.offsetWidth||0,r=d-e.clientX;n+=t.placement===`right`?r:-r,n=A(n),k(n),d=e.clientX}}function P(){_.value&&(d=0,_.value=!1,document.body.style.cursor=f,document.body.removeEventListener(`mousemove`,N),document.body.removeEventListener(`mouseup`,P),document.body.removeEventListener(`mouseleave`,w))}b(()=>{t.show&&(r.value=!0)}),e(()=>t.show,e=>{e||P()}),g(()=>{P()});let F=C(()=>{let{show:e}=t,n=[[L,e]];return t.showMask||n.push([I,t.onClickoutside,void 0,{capture:!0}]),n});function R(){var e;r.value=!1,(e=t.onAfterLeave)==null||e.call(t)}return z(C(()=>t.blockScroll&&r.value)),c(i,s),c(h,null),c(o,null),{bodyRef:s,rtlEnabled:S,mergedClsPrefix:l.mergedClsPrefixRef,isMounted:l.isMountedRef,mergedTheme:l.mergedThemeRef,displayed:r,transitionName:C(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[t.placement]),handleAfterLeave:R,bodyDirectives:F,handleMousedownResizeTrigger:T,handleMouseenterResizeTrigger:E,handleMouseleaveResizeTrigger:D,isDragging:_,isHoverOnResizeTrigger:m}},render(){let{$slots:e,mergedClsPrefix:n}=this;return this.displayDirective===`show`||this.displayed||this.show?E(f(`div`,{role:`none`},f(T,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>f(R,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>E(f(`div`,t(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${n}-drawer`,this.rtlEnabled&&`${n}-drawer--rtl`,`${n}-drawer--${this.placement}-placement`,this.isDragging&&`${n}-drawer--unselectable`,this.nativeScrollbar&&`${n}-drawer--native-scrollbar`]}),[this.resizable?f(`div`,{class:[`${n}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${n}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?f(`div`,{class:[`${n}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:`none`},e):f(A,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${n}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),e)]),this.bodyDirectives)})})),[[L,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:W,cubicBezierEaseOut:G}=m;function K({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${W}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${G}`}),x(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:q,cubicBezierEaseOut:J}=m;function Y({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${q}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${J}`}),x(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:X,cubicBezierEaseOut:Z}=m;function Q({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${X}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${Z}`}),x(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:$,cubicBezierEaseOut:ee}=m;function te({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${$}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${ee}`}),x(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var ne=x([P(`drawer`,`
 word-break: break-word;
 line-height: var(--n-line-height);
 position: absolute;
 pointer-events: all;
 box-shadow: var(--n-box-shadow);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 background-color: var(--n-color);
 color: var(--n-text-color);
 box-sizing: border-box;
 `,[Q(),Y(),te(),K(),k(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),k(`native-scrollbar`,[P(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),D(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[k(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),P(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),P(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[k(`native-scrollbar`,[P(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),P(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),P(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),P(`drawer-header`,`
 font-weight: var(--n-title-font-weight);
 line-height: 1;
 font-size: var(--n-title-font-size);
 color: var(--n-title-text-color);
 padding: var(--n-header-padding);
 transition: border .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-divider-color);
 border-bottom: var(--n-header-border-bottom);
 display: flex;
 justify-content: space-between;
 align-items: center;
 `,[D(`main`,`
 flex: 1;
 `),D(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),P(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),k(`right-placement`,`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[D(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),k(`left-placement`,`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[D(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),k(`top-placement`,`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[D(`resize-trigger`,`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),k(`bottom-placement`,`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[D(`resize-trigger`,`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),x(`body`,[x(`>`,[P(`drawer-container`,`
 position: fixed;
 `)])]),P(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[x(`> *`,`
 pointer-events: all;
 `)]),P(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[k(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),M({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),re=l({name:`Drawer`,inheritAttrs:!1,props:Object.assign(Object.assign({},r.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),setup(e){let{mergedClsPrefixRef:t,namespaceRef:i,inlineThemeDisabled:a}=n(e),o=d(),l=r(`Drawer`,`-drawer`,ne,H,e,t),f=j(e.defaultWidth),m=j(e.defaultHeight),h=s(w(e,`width`),f),g=s(w(e,`height`),m),_=C(()=>{let{placement:t}=e;return t===`top`||t===`bottom`?``:N(h.value)}),y=C(()=>{let{placement:t}=e;return t===`left`||t===`right`?``:N(g.value)}),b=t=>{let{onUpdateWidth:n,"onUpdate:width":r}=e;n&&S(n,t),r&&S(r,t),f.value=t},x=t=>{let{onUpdateHeight:n,"onUpdate:width":r}=e;n&&S(n,t),r&&S(r,t),m.value=t},T=C(()=>[{width:_.value,height:y.value},e.drawerStyle||``]);function E(t){let{onMaskClick:n,maskClosable:r}=e;r&&A(!1),n&&n(t)}function D(e){E(e)}let O=p();function k(t){var n;(n=e.onEsc)==null||n.call(e),e.show&&e.closeOnEsc&&v(t)&&(O.value||A(!1))}function A(t){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=e;r&&S(r,t),i&&S(i,t),n&&!t&&S(n,t)}c(u,{isMountedRef:o,mergedThemeRef:l,mergedClsPrefixRef:t,doUpdateShow:A,doUpdateHeight:x,doUpdateWidth:b});let M=C(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:a,lineHeight:o,headerPadding:s,footerPadding:c,borderRadius:u,bodyPadding:d,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=l.value;return{"--n-line-height":o,"--n-color":r,"--n-border-radius":u,"--n-text-color":i,"--n-box-shadow":a,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":s,"--n-body-padding":d,"--n-footer-padding":c,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),P=a?V(`drawer`,void 0,M,e):void 0;return{mergedClsPrefix:t,namespace:i,mergedBodyStyle:T,handleOutsideClick:D,handleMaskClick:E,handleEsc:k,mergedTheme:l,cssVars:a?void 0:M,themeClass:P?.themeClass,onRender:P?.onRender,isMounted:o}},render(){let{mergedClsPrefix:e}=this;return f(O,{to:this.to,show:this.show},{default:()=>{var t;return(t=this.onRender)==null||t.call(this),E(f(`div`,{class:[`${e}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:`none`},this.showMask?f(R,{name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?f(`div`,{"aria-hidden":!0,class:[`${e}-drawer-mask`,this.showMask===`transparent`&&`${e}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,f(U,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[_,{zIndex:this.zIndex,enabled:this.show}]])}})}}),ie=l({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let e=B(u,null);e||F(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:t}=e;function n(){t(!1)}return{handleCloseClick:n,mergedTheme:e.mergedThemeRef,mergedClsPrefix:e.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:t,nativeScrollbar:n,mergedTheme:r,bodyClass:i,bodyStyle:a,bodyContentClass:o,bodyContentStyle:s,headerClass:c,headerStyle:l,footerClass:u,footerStyle:d,scrollbarProps:p,closable:m,$slots:h}=this;return f(`div`,{role:`none`,class:[`${t}-drawer-content`,n&&`${t}-drawer-content--native-scrollbar`]},h.header||e||m?f(`div`,{class:[`${t}-drawer-header`,c],style:l,role:`none`},f(`div`,{class:`${t}-drawer-header__main`,role:`heading`,"aria-level":`1`},h.header===void 0?e:h.header()),m&&f(y,{onClick:this.handleCloseClick,clsPrefix:t,class:`${t}-drawer-header__close`,absolute:!0})):null,n?f(`div`,{class:[`${t}-drawer-body`,i],style:a,role:`none`},f(`div`,{class:[`${t}-drawer-body-content-wrapper`,o],style:s,role:`none`},h)):f(A,Object.assign({themeOverrides:r.peerOverrides.Scrollbar,theme:r.peers.Scrollbar},p,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,o],contentStyle:s}),h),h.footer?f(`div`,{class:[`${t}-drawer-footer`,u],style:d,role:`none`},h.footer()):null)}});export{re as n,ie as t};