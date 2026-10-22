import{An as e,Ar as t,Dn as n,Er as r,Hi as i,Mr as a,Ni as o,Oi as s,Qi as c,Rr as l,Sr as u,Ti as d,Tr as f,Ui as p,Un as m,Vi as h,Wi as g,Wn as _,_r as v,ai as y,ar as b,da as x,dn as S,dr as C,hn as w,ia as T,ii as E,ir as D,jr as O,kn as k,kr as A,lr as j,oa as M,oi as N,pa as P,pr as F,qi as I,ra as L,si as R,sn as z,xr as B,zr as V}from"./router-DhxDhlAH.js";import{r as H}from"./index-DHsLjxY8.js";var U=h({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(n){let r=x(!!n.show),i=x(null),l=p(a),d=0,m=``,h=null,g=x(!1),v=x(!1),y=o(()=>n.placement===`top`||n.placement===`bottom`),{mergedClsPrefixRef:b,mergedRtlRef:S}=_(n),C=e(`Drawer`,S,b),w=R,E=e=>{v.value=!0,d=y.value?e.clientY:e.clientX,m=document.body.style.cursor,document.body.style.cursor=y.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,F),document.body.addEventListener(`mouseleave`,w),document.body.addEventListener(`mouseup`,R)},D=()=>{h!==null&&(window.clearTimeout(h),h=null),v.value?g.value=!0:h=window.setTimeout(()=>{g.value=!0},300)},k=()=>{h!==null&&(window.clearTimeout(h),h=null),g.value=!1},{doUpdateHeight:j,doUpdateWidth:M}=l,N=e=>{let{maxWidth:t}=n;if(t&&e>t)return t;let{minWidth:r}=n;return r&&e<r?r:e},P=e=>{let{maxHeight:t}=n;if(t&&e>t)return t;let{minHeight:r}=n;return r&&e<r?r:e};function F(e){if(v.value)if(y.value){let t=i.value?.offsetHeight||0,r=d-e.clientY;t+=n.placement===`bottom`?r:-r,t=P(t),j(t),d=e.clientY}else{let t=i.value?.offsetWidth||0,r=d-e.clientX;t+=n.placement===`right`?r:-r,t=N(t),M(t),d=e.clientX}}function R(){v.value&&(d=0,v.value=!1,document.body.style.cursor=m,document.body.removeEventListener(`mousemove`,F),document.body.removeEventListener(`mouseup`,R),document.body.removeEventListener(`mouseleave`,w))}T(()=>{n.show&&(r.value=!0)}),L(()=>n.show,e=>{e||R()}),I(()=>{R()});let z=o(()=>{let{show:e}=n,t=[[s,e]];return n.showMask||t.push([u,n.onClickoutside,void 0,{capture:!0}]),t});function B(){var e;r.value=!1,(e=n.onAfterLeave)==null||e.call(n)}return f(o(()=>n.blockScroll&&r.value)),c(O,i),c(A,null),c(t,null),{bodyRef:i,rtlEnabled:C,mergedClsPrefix:l.mergedClsPrefixRef,isMounted:l.isMountedRef,mergedTheme:l.mergedThemeRef,displayed:r,transitionName:o(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[n.placement]),handleAfterLeave:B,bodyDirectives:z,handleMousedownResizeTrigger:E,handleMouseenterResizeTrigger:D,handleMouseleaveResizeTrigger:k,isDragging:v,isHoverOnResizeTrigger:g}},render(){let{$slots:e,mergedClsPrefix:t}=this;return this.displayDirective===`show`||this.displayed||this.show?M(i(`div`,{role:`none`},i(F,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>i(d,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>M(i(`div`,g(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${t}-drawer`,this.rtlEnabled&&`${t}-drawer--rtl`,`${t}-drawer--${this.placement}-placement`,this.isDragging&&`${t}-drawer--unselectable`,this.nativeScrollbar&&`${t}-drawer--native-scrollbar`]}),[this.resizable?i(`div`,{class:[`${t}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${t}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?i(`div`,{class:[`${t}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:`none`},e):i(z,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${t}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),e)]),this.bodyDirectives)})})),[[s,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:W,cubicBezierEaseOut:G}=k;function K({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${W}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${G}`}),E(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:q,cubicBezierEaseOut:J}=k;function Y({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${q}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${J}`}),E(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:X,cubicBezierEaseOut:Z}=k;function Q({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${X}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${Z}`}),E(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:$,cubicBezierEaseOut:ee}=k;function te({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${$}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${ee}`}),E(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var ne=E([y(`drawer`,`
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
 `,[Q(),Y(),te(),K(),R(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),R(`native-scrollbar`,[y(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),N(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[R(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),y(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),y(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[R(`native-scrollbar`,[y(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),y(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),y(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),y(`drawer-header`,`
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
 `,[N(`main`,`
 flex: 1;
 `),N(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),y(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),R(`right-placement`,`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),R(`left-placement`,`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),R(`top-placement`,`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),R(`bottom-placement`,`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),E(`body`,[E(`>`,[y(`drawer-container`,`
 position: fixed;
 `)])]),y(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[E(`> *`,`
 pointer-events: all;
 `)]),y(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[R(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),S({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),re=h({name:`Drawer`,inheritAttrs:!1,props:Object.assign(Object.assign({},n.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),setup(e){let{mergedClsPrefixRef:t,namespaceRef:i,inlineThemeDisabled:s}=_(e),u=l(),d=n(`Drawer`,`-drawer`,ne,H,e,t),f=x(e.defaultWidth),p=x(e.defaultHeight),h=V(P(e,`width`),f),g=V(P(e,`height`),p),v=o(()=>{let{placement:t}=e;return t===`top`||t===`bottom`?``:C(h.value)}),y=o(()=>{let{placement:t}=e;return t===`left`||t===`right`?``:C(g.value)}),b=t=>{let{onUpdateWidth:n,"onUpdate:width":r}=e;n&&D(n,t),r&&D(r,t),f.value=t},S=t=>{let{onUpdateHeight:n,"onUpdate:width":r}=e;n&&D(n,t),r&&D(r,t),p.value=t},w=o(()=>[{width:v.value,height:y.value},e.drawerStyle||``]);function T(t){let{onMaskClick:n,maskClosable:r}=e;r&&A(!1),n&&n(t)}function E(e){T(e)}let O=r();function k(t){var n;(n=e.onEsc)==null||n.call(e),e.show&&e.closeOnEsc&&j(t)&&(O.value||A(!1))}function A(t){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=e;r&&D(r,t),i&&D(i,t),n&&!t&&D(n,t)}c(a,{isMountedRef:u,mergedThemeRef:d,mergedClsPrefixRef:t,doUpdateShow:A,doUpdateHeight:S,doUpdateWidth:b});let M=o(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:a,lineHeight:o,headerPadding:s,footerPadding:c,borderRadius:l,bodyPadding:u,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=d.value;return{"--n-line-height":o,"--n-color":r,"--n-border-radius":l,"--n-text-color":i,"--n-box-shadow":a,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":s,"--n-body-padding":u,"--n-footer-padding":c,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),N=s?m(`drawer`,void 0,M,e):void 0;return{mergedClsPrefix:t,namespace:i,mergedBodyStyle:w,handleOutsideClick:E,handleMaskClick:T,handleEsc:k,mergedTheme:d,cssVars:s?void 0:M,themeClass:N?.themeClass,onRender:N?.onRender,isMounted:u}},render(){let{mergedClsPrefix:e}=this;return i(v,{to:this.to,show:this.show},{default:()=>{var t;return(t=this.onRender)==null||t.call(this),M(i(`div`,{class:[`${e}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:`none`},this.showMask?i(d,{name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?i(`div`,{"aria-hidden":!0,class:[`${e}-drawer-mask`,this.showMask===`transparent`&&`${e}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,i(U,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[B,{zIndex:this.zIndex,enabled:this.show}]])}})}}),ie=h({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let e=p(a,null);e||b(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:t}=e;function n(){t(!1)}return{handleCloseClick:n,mergedTheme:e.mergedThemeRef,mergedClsPrefix:e.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:t,nativeScrollbar:n,mergedTheme:r,bodyClass:a,bodyStyle:o,bodyContentClass:s,bodyContentStyle:c,headerClass:l,headerStyle:u,footerClass:d,footerStyle:f,scrollbarProps:p,closable:m,$slots:h}=this;return i(`div`,{role:`none`,class:[`${t}-drawer-content`,n&&`${t}-drawer-content--native-scrollbar`]},h.header||e||m?i(`div`,{class:[`${t}-drawer-header`,l],style:u,role:`none`},i(`div`,{class:`${t}-drawer-header__main`,role:`heading`,"aria-level":`1`},h.header===void 0?e:h.header()),m&&i(w,{onClick:this.handleCloseClick,clsPrefix:t,class:`${t}-drawer-header__close`,absolute:!0})):null,n?i(`div`,{class:[`${t}-drawer-body`,a],style:o,role:`none`},i(`div`,{class:[`${t}-drawer-body-content-wrapper`,s],style:c,role:`none`},h)):i(z,Object.assign({themeOverrides:r.peerOverrides.Scrollbar,theme:r.peers.Scrollbar},p,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,s],contentStyle:c}),h),h.footer?i(`div`,{class:[`${t}-drawer-footer`,d],style:f,role:`none`},h.footer()):null)}});export{re as n,ie as t};