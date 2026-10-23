import{Ar as e,Er as t,Fr as n,Gi as r,Hr as i,Ir as a,Ji as o,Jn as s,Ki as c,Li as l,Mi as u,Nn as d,Nr as f,Pn as p,Pr as m,Tr as h,Ur as g,Zi as _,_a as v,_r as y,ca as b,ci as x,cr as S,di as C,dn as w,ha as T,hn as E,hr as D,jn as O,ki as k,kr as A,li as j,lr as M,na as N,pr as P,qi as F,qn as I,sa as L,ua as R,ui as z,xr as B,yn as V}from"./router-DnIJc1nz.js";import{r as H}from"./index-CQ6w_Oy4.js";var U=r({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(e){let r=T(!!e.show),i=T(null),o=F(a),c=0,d=``,h=null,g=T(!1),v=T(!1),y=l(()=>e.placement===`top`||e.placement===`bottom`),{mergedClsPrefixRef:x,mergedRtlRef:S}=s(e),C=p(`Drawer`,S,x),w=R,E=e=>{v.value=!0,c=y.value?e.clientY:e.clientX,d=document.body.style.cursor,document.body.style.cursor=y.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,I),document.body.addEventListener(`mouseleave`,w),document.body.addEventListener(`mouseup`,R)},D=()=>{h!==null&&(window.clearTimeout(h),h=null),v.value?g.value=!0:h=window.setTimeout(()=>{g.value=!0},300)},O=()=>{h!==null&&(window.clearTimeout(h),h=null),g.value=!1},{doUpdateHeight:k,doUpdateWidth:j}=o,M=t=>{let{maxWidth:n}=e;if(n&&t>n)return n;let{minWidth:r}=e;return r&&t<r?r:t},P=t=>{let{maxHeight:n}=e;if(n&&t>n)return n;let{minHeight:r}=e;return r&&t<r?r:t};function I(t){if(v.value)if(y.value){let n=i.value?.offsetHeight||0,r=c-t.clientY;n+=e.placement===`bottom`?r:-r,n=P(n),k(n),c=t.clientY}else{let n=i.value?.offsetWidth||0,r=c-t.clientX;n+=e.placement===`right`?r:-r,n=M(n),j(n),c=t.clientX}}function R(){v.value&&(c=0,v.value=!1,document.body.style.cursor=d,document.body.removeEventListener(`mousemove`,I),document.body.removeEventListener(`mouseup`,R),document.body.removeEventListener(`mouseleave`,w))}b(()=>{e.show&&(r.value=!0)}),L(()=>e.show,e=>{e||R()}),_(()=>{R()});let z=l(()=>{let{show:n}=e,r=[[u,n]];return e.showMask||r.push([t,e.onClickoutside,void 0,{capture:!0}]),r});function B(){var t;r.value=!1,(t=e.onAfterLeave)==null||t.call(e)}return A(l(()=>e.blockScroll&&r.value)),N(n,i),N(f,null),N(m,null),{bodyRef:i,rtlEnabled:C,mergedClsPrefix:o.mergedClsPrefixRef,isMounted:o.isMountedRef,mergedTheme:o.mergedThemeRef,displayed:r,transitionName:l(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[e.placement]),handleAfterLeave:B,bodyDirectives:z,handleMousedownResizeTrigger:E,handleMouseenterResizeTrigger:D,handleMouseleaveResizeTrigger:O,isDragging:v,isHoverOnResizeTrigger:g}},render(){let{$slots:e,mergedClsPrefix:t}=this;return this.displayDirective===`show`||this.displayed||this.show?R(c(`div`,{role:`none`},c(y,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>c(k,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>R(c(`div`,o(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${t}-drawer`,this.rtlEnabled&&`${t}-drawer--rtl`,`${t}-drawer--${this.placement}-placement`,this.isDragging&&`${t}-drawer--unselectable`,this.nativeScrollbar&&`${t}-drawer--native-scrollbar`]}),[this.resizable?c(`div`,{class:[`${t}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${t}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?c(`div`,{class:[`${t}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:`none`},e):c(w,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${t}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),e)]),this.bodyDirectives)})})),[[u,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:W,cubicBezierEaseOut:G}=d;function K({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${W}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${G}`}),x(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:q,cubicBezierEaseOut:J}=d;function Y({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${q}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${J}`}),x(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:X,cubicBezierEaseOut:Z}=d;function Q({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${X}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${Z}`}),x(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:$,cubicBezierEaseOut:ee}=d;function te({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[x(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${$}`}),x(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${ee}`}),x(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),x(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),x(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),x(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var ne=x([j(`drawer`,`
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
 `,[Q(),Y(),te(),K(),C(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),C(`native-scrollbar`,[j(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),z(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[C(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),j(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),j(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[C(`native-scrollbar`,[j(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),j(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),j(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),j(`drawer-header`,`
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
 `,[z(`main`,`
 flex: 1;
 `),z(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),j(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),C(`right-placement`,`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[z(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),C(`left-placement`,`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[z(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),C(`top-placement`,`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[z(`resize-trigger`,`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),C(`bottom-placement`,`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[z(`resize-trigger`,`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),x(`body`,[x(`>`,[j(`drawer-container`,`
 position: fixed;
 `)])]),j(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[x(`> *`,`
 pointer-events: all;
 `)]),j(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[C(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),E({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),re=r({name:`Drawer`,inheritAttrs:!1,props:Object.assign(Object.assign({},O.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),setup(t){let{mergedClsPrefixRef:n,namespaceRef:r,inlineThemeDisabled:o}=s(t),c=i(),u=O(`Drawer`,`-drawer`,ne,H,t,n),d=T(t.defaultWidth),f=T(t.defaultHeight),p=g(v(t,`width`),d),m=g(v(t,`height`),f),h=l(()=>{let{placement:e}=t;return e===`top`||e===`bottom`?``:D(p.value)}),_=l(()=>{let{placement:e}=t;return e===`left`||e===`right`?``:D(m.value)}),y=e=>{let{onUpdateWidth:n,"onUpdate:width":r}=t;n&&S(n,e),r&&S(r,e),d.value=e},b=e=>{let{onUpdateHeight:n,"onUpdate:width":r}=t;n&&S(n,e),r&&S(r,e),f.value=e},x=l(()=>[{width:h.value,height:_.value},t.drawerStyle||``]);function C(e){let{onMaskClick:n,maskClosable:r}=t;r&&A(!1),n&&n(e)}function w(e){C(e)}let E=e();function k(e){var n;(n=t.onEsc)==null||n.call(t),t.show&&t.closeOnEsc&&P(e)&&(E.value||A(!1))}function A(e){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=t;r&&S(r,e),i&&S(i,e),n&&!e&&S(n,e)}N(a,{isMountedRef:c,mergedThemeRef:u,mergedClsPrefixRef:n,doUpdateShow:A,doUpdateHeight:b,doUpdateWidth:y});let j=l(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:a,lineHeight:o,headerPadding:s,footerPadding:c,borderRadius:l,bodyPadding:d,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=u.value;return{"--n-line-height":o,"--n-color":r,"--n-border-radius":l,"--n-text-color":i,"--n-box-shadow":a,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":s,"--n-body-padding":d,"--n-footer-padding":c,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),M=o?I(`drawer`,void 0,j,t):void 0;return{mergedClsPrefix:n,namespace:r,mergedBodyStyle:x,handleOutsideClick:w,handleMaskClick:C,handleEsc:k,mergedTheme:u,cssVars:o?void 0:j,themeClass:M?.themeClass,onRender:M?.onRender,isMounted:c}},render(){let{mergedClsPrefix:e}=this;return c(B,{to:this.to,show:this.show},{default:()=>{var t;return(t=this.onRender)==null||t.call(this),R(c(`div`,{class:[`${e}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:`none`},this.showMask?c(k,{name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?c(`div`,{"aria-hidden":!0,class:[`${e}-drawer-mask`,this.showMask===`transparent`&&`${e}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,c(U,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[h,{zIndex:this.zIndex,enabled:this.show}]])}})}}),ie=r({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let e=F(a,null);e||M(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:t}=e;function n(){t(!1)}return{handleCloseClick:n,mergedTheme:e.mergedThemeRef,mergedClsPrefix:e.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:t,nativeScrollbar:n,mergedTheme:r,bodyClass:i,bodyStyle:a,bodyContentClass:o,bodyContentStyle:s,headerClass:l,headerStyle:u,footerClass:d,footerStyle:f,scrollbarProps:p,closable:m,$slots:h}=this;return c(`div`,{role:`none`,class:[`${t}-drawer-content`,n&&`${t}-drawer-content--native-scrollbar`]},h.header||e||m?c(`div`,{class:[`${t}-drawer-header`,l],style:u,role:`none`},c(`div`,{class:`${t}-drawer-header__main`,role:`heading`,"aria-level":`1`},h.header===void 0?e:h.header()),m&&c(V,{onClick:this.handleCloseClick,clsPrefix:t,class:`${t}-drawer-header__close`,absolute:!0})):null,n?c(`div`,{class:[`${t}-drawer-body`,i],style:a,role:`none`},c(`div`,{class:[`${t}-drawer-body-content-wrapper`,o],style:s,role:`none`},h)):c(w,Object.assign({themeOverrides:r.peerOverrides.Scrollbar,theme:r.peers.Scrollbar},p,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,o],contentStyle:s}),h),h.footer?c(`div`,{class:[`${t}-drawer-footer`,d],style:f,role:`none`},h.footer()):null)}});export{re as n,ie as t};