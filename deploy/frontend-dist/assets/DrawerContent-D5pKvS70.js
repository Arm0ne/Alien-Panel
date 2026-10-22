import{Ar as e,Bi as t,Ci as n,Cr as r,Dn as i,Dr as a,Ei as o,Gi as s,Hi as c,Hn as l,Ir as u,Lr as d,On as f,Or as p,Tn as m,Vi as h,Vn as g,Xi as _,ai as v,an as y,br as b,da as x,dr as S,hr as C,ia as w,ii as T,ji as E,kr as D,la as O,ln as k,lr as A,na as j,ni as M,nr as N,pn as P,ri as F,rr as I,sr as L,ta as R,wr as z,yr as B,zi as V}from"./router-DMnLqP4l.js";import{r as H}from"./index-c6iM4-gc.js";var U=V({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(t){let n=O(!!t.show),i=O(null),c=h(e),u=0,d=``,m=null,g=O(!1),v=O(!1),y=E(()=>t.placement===`top`||t.placement===`bottom`),{mergedClsPrefixRef:x,mergedRtlRef:S}=l(t),C=f(`Drawer`,S,x),w=L,T=e=>{v.value=!0,u=y.value?e.clientY:e.clientX,d=document.body.style.cursor,document.body.style.cursor=y.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,I),document.body.addEventListener(`mouseleave`,w),document.body.addEventListener(`mouseup`,L)},k=()=>{m!==null&&(window.clearTimeout(m),m=null),v.value?g.value=!0:m=window.setTimeout(()=>{g.value=!0},300)},A=()=>{m!==null&&(window.clearTimeout(m),m=null),g.value=!1},{doUpdateHeight:M,doUpdateWidth:N}=c,P=e=>{let{maxWidth:n}=t;if(n&&e>n)return n;let{minWidth:r}=t;return r&&e<r?r:e},F=e=>{let{maxHeight:n}=t;if(n&&e>n)return n;let{minHeight:r}=t;return r&&e<r?r:e};function I(e){if(v.value)if(y.value){let n=i.value?.offsetHeight||0,r=u-e.clientY;n+=t.placement===`bottom`?r:-r,n=F(n),M(n),u=e.clientY}else{let n=i.value?.offsetWidth||0,r=u-e.clientX;n+=t.placement===`right`?r:-r,n=P(n),N(n),u=e.clientX}}function L(){v.value&&(u=0,v.value=!1,document.body.style.cursor=d,document.body.removeEventListener(`mousemove`,I),document.body.removeEventListener(`mouseup`,L),document.body.removeEventListener(`mouseleave`,w))}j(()=>{t.show&&(n.value=!0)}),R(()=>t.show,e=>{e||L()}),s(()=>{L()});let z=E(()=>{let{show:e}=t,n=[[o,e]];return t.showMask||n.push([b,t.onClickoutside,void 0,{capture:!0}]),n});function B(){var e;n.value=!1,(e=t.onAfterLeave)==null||e.call(t)}return r(E(()=>t.blockScroll&&n.value)),_(D,i),_(a,null),_(p,null),{bodyRef:i,rtlEnabled:C,mergedClsPrefix:c.mergedClsPrefixRef,isMounted:c.isMountedRef,mergedTheme:c.mergedThemeRef,displayed:n,transitionName:E(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[t.placement]),handleAfterLeave:B,bodyDirectives:z,handleMousedownResizeTrigger:T,handleMouseenterResizeTrigger:k,handleMouseleaveResizeTrigger:A,isDragging:v,isHoverOnResizeTrigger:g}},render(){let{$slots:e,mergedClsPrefix:r}=this;return this.displayDirective===`show`||this.displayed||this.show?w(t(`div`,{role:`none`},t(S,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>t(n,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>w(t(`div`,c(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${r}-drawer`,this.rtlEnabled&&`${r}-drawer--rtl`,`${r}-drawer--${this.placement}-placement`,this.isDragging&&`${r}-drawer--unselectable`,this.nativeScrollbar&&`${r}-drawer--native-scrollbar`]}),[this.resizable?t(`div`,{class:[`${r}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${r}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?t(`div`,{class:[`${r}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:`none`},e):t(y,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${r}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),e)]),this.bodyDirectives)})})),[[o,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:W,cubicBezierEaseOut:G}=i;function K({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[M(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${W}`}),M(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${G}`}),M(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),M(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),M(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),M(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:q,cubicBezierEaseOut:J}=i;function Y({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[M(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${q}`}),M(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${J}`}),M(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),M(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),M(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),M(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:X,cubicBezierEaseOut:Z}=i;function Q({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[M(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${X}`}),M(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${Z}`}),M(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),M(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),M(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),M(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:$,cubicBezierEaseOut:ee}=i;function te({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[M(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${$}`}),M(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${ee}`}),M(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),M(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),M(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),M(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var ne=M([F(`drawer`,`
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
 `,[Q(),Y(),te(),K(),v(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),v(`native-scrollbar`,[F(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),T(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[v(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),F(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),F(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[v(`native-scrollbar`,[F(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),F(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),F(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),F(`drawer-header`,`
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
 `,[T(`main`,`
 flex: 1;
 `),T(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),F(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),v(`right-placement`,`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[T(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),v(`left-placement`,`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[T(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),v(`top-placement`,`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[T(`resize-trigger`,`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),v(`bottom-placement`,`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[T(`resize-trigger`,`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),M(`body`,[M(`>`,[F(`drawer-container`,`
 position: fixed;
 `)])]),F(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[M(`> *`,`
 pointer-events: all;
 `)]),F(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[v(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),k({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),re=V({name:`Drawer`,inheritAttrs:!1,props:Object.assign(Object.assign({},m.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),setup(t){let{mergedClsPrefixRef:n,namespaceRef:r,inlineThemeDisabled:i}=l(t),a=u(),o=m(`Drawer`,`-drawer`,ne,H,t,n),s=O(t.defaultWidth),c=O(t.defaultHeight),f=d(x(t,`width`),s),p=d(x(t,`height`),c),h=E(()=>{let{placement:e}=t;return e===`top`||e===`bottom`?``:A(f.value)}),v=E(()=>{let{placement:e}=t;return e===`left`||e===`right`?``:A(p.value)}),y=e=>{let{onUpdateWidth:n,"onUpdate:width":r}=t;n&&N(n,e),r&&N(r,e),s.value=e},b=e=>{let{onUpdateHeight:n,"onUpdate:width":r}=t;n&&N(n,e),r&&N(r,e),c.value=e},S=E(()=>[{width:h.value,height:v.value},t.drawerStyle||``]);function C(e){let{onMaskClick:n,maskClosable:r}=t;r&&k(!1),n&&n(e)}function w(e){C(e)}let T=z();function D(e){var n;(n=t.onEsc)==null||n.call(t),t.show&&t.closeOnEsc&&L(e)&&(T.value||k(!1))}function k(e){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=t;r&&N(r,e),i&&N(i,e),n&&!e&&N(n,e)}_(e,{isMountedRef:a,mergedThemeRef:o,mergedClsPrefixRef:n,doUpdateShow:k,doUpdateHeight:b,doUpdateWidth:y});let j=E(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:a,lineHeight:s,headerPadding:c,footerPadding:l,borderRadius:u,bodyPadding:d,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=o.value;return{"--n-line-height":s,"--n-color":r,"--n-border-radius":u,"--n-text-color":i,"--n-box-shadow":a,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":c,"--n-body-padding":d,"--n-footer-padding":l,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),M=i?g(`drawer`,void 0,j,t):void 0;return{mergedClsPrefix:n,namespace:r,mergedBodyStyle:S,handleOutsideClick:w,handleMaskClick:C,handleEsc:D,mergedTheme:o,cssVars:i?void 0:j,themeClass:M?.themeClass,onRender:M?.onRender,isMounted:a}},render(){let{mergedClsPrefix:e}=this;return t(C,{to:this.to,show:this.show},{default:()=>{var r;return(r=this.onRender)==null||r.call(this),w(t(`div`,{class:[`${e}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:`none`},this.showMask?t(n,{name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?t(`div`,{"aria-hidden":!0,class:[`${e}-drawer-mask`,this.showMask===`transparent`&&`${e}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,t(U,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[B,{zIndex:this.zIndex,enabled:this.show}]])}})}}),ie=V({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let t=h(e,null);t||I(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:n}=t;function r(){n(!1)}return{handleCloseClick:r,mergedTheme:t.mergedThemeRef,mergedClsPrefix:t.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:n,nativeScrollbar:r,mergedTheme:i,bodyClass:a,bodyStyle:o,bodyContentClass:s,bodyContentStyle:c,headerClass:l,headerStyle:u,footerClass:d,footerStyle:f,scrollbarProps:p,closable:m,$slots:h}=this;return t(`div`,{role:`none`,class:[`${n}-drawer-content`,r&&`${n}-drawer-content--native-scrollbar`]},h.header||e||m?t(`div`,{class:[`${n}-drawer-header`,l],style:u,role:`none`},t(`div`,{class:`${n}-drawer-header__main`,role:`heading`,"aria-level":`1`},h.header===void 0?e:h.header()),m&&t(P,{onClick:this.handleCloseClick,clsPrefix:n,class:`${n}-drawer-header__close`,absolute:!0})):null,r?t(`div`,{class:[`${n}-drawer-body`,a],style:o,role:`none`},t(`div`,{class:[`${n}-drawer-body-content-wrapper`,s],style:c,role:`none`},h)):t(y,Object.assign({themeOverrides:i.peerOverrides.Scrollbar,theme:i.peers.Scrollbar},p,{class:`${n}-drawer-body`,contentClass:[`${n}-drawer-body-content-wrapper`,s],contentStyle:c}),h),h.footer?t(`div`,{class:[`${n}-drawer-footer`,d],style:f,role:`none`},h.footer()):null)}});export{re as n,ie as t};