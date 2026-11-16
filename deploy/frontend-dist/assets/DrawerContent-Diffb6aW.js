import{Ai as e,Ar as t,Dr as n,Er as r,Fn as i,Fr as a,Ir as o,Ji as s,Jn as c,Ki as l,Lr as u,Mn as d,Ni as f,Pn as p,Pr as m,Qi as h,Ri as g,Sr as _,Ur as v,Wr as y,Yi as b,Yn as x,bn as S,ca as C,da as w,di as T,fi as E,fn as D,ga as O,gn as k,gr as A,jr as j,la as M,li as N,lr as P,mr as F,qi as I,ra as L,ui as R,ur as z,va as B,vr as V}from"./router-CYJ0icnP.js";import{r as H}from"./index-R29AZ2H8.js";var U=l({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(e){let r=O(!!e.show),c=O(null),l=s(u),d=0,p=``,_=null,v=O(!1),y=O(!1),b=g(()=>e.placement===`top`||e.placement===`bottom`),{mergedClsPrefixRef:S,mergedRtlRef:w}=x(e),T=i(`Drawer`,w,S),E=R,D=e=>{y.value=!0,d=b.value?e.clientY:e.clientX,p=document.body.style.cursor,document.body.style.cursor=b.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,I),document.body.addEventListener(`mouseleave`,E),document.body.addEventListener(`mouseup`,R)},k=()=>{_!==null&&(window.clearTimeout(_),_=null),y.value?v.value=!0:_=window.setTimeout(()=>{v.value=!0},300)},A=()=>{_!==null&&(window.clearTimeout(_),_=null),v.value=!1},{doUpdateHeight:j,doUpdateWidth:N}=l,P=t=>{let{maxWidth:n}=e;if(n&&t>n)return n;let{minWidth:r}=e;return r&&t<r?r:t},F=t=>{let{maxHeight:n}=e;if(n&&t>n)return n;let{minHeight:r}=e;return r&&t<r?r:t};function I(t){if(y.value)if(b.value){let n=c.value?.offsetHeight||0,r=d-t.clientY;n+=e.placement===`bottom`?r:-r,n=F(n),j(n),d=t.clientY}else{let n=c.value?.offsetWidth||0,r=d-t.clientX;n+=e.placement===`right`?r:-r,n=P(n),N(n),d=t.clientX}}function R(){y.value&&(d=0,y.value=!1,document.body.style.cursor=p,document.body.removeEventListener(`mousemove`,I),document.body.removeEventListener(`mouseup`,R),document.body.removeEventListener(`mouseleave`,E))}M(()=>{e.show&&(r.value=!0)}),C(()=>e.show,e=>{e||R()}),h(()=>{R()});let z=g(()=>{let{show:t}=e,r=[[f,t]];return e.showMask||r.push([n,e.onClickoutside,void 0,{capture:!0}]),r});function B(){var t;r.value=!1,(t=e.onAfterLeave)==null||t.call(e)}return t(g(()=>e.blockScroll&&r.value)),L(o,c),L(m,null),L(a,null),{bodyRef:c,rtlEnabled:T,mergedClsPrefix:l.mergedClsPrefixRef,isMounted:l.isMountedRef,mergedTheme:l.mergedThemeRef,displayed:r,transitionName:g(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[e.placement]),handleAfterLeave:B,bodyDirectives:z,handleMousedownResizeTrigger:D,handleMouseenterResizeTrigger:k,handleMouseleaveResizeTrigger:A,isDragging:y,isHoverOnResizeTrigger:v}},render(){let{$slots:t,mergedClsPrefix:n}=this;return this.displayDirective===`show`||this.displayed||this.show?w(I(`div`,{role:`none`},I(V,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>I(e,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>w(I(`div`,b(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${n}-drawer`,this.rtlEnabled&&`${n}-drawer--rtl`,`${n}-drawer--${this.placement}-placement`,this.isDragging&&`${n}-drawer--unselectable`,this.nativeScrollbar&&`${n}-drawer--native-scrollbar`]}),[this.resizable?I(`div`,{class:[`${n}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${n}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?I(`div`,{class:[`${n}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:`none`},t):I(D,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${n}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),t)]),this.bodyDirectives)})})),[[f,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:W,cubicBezierEaseOut:G}=p;function K({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[N(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${W}`}),N(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${G}`}),N(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),N(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),N(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),N(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:q,cubicBezierEaseOut:J}=p;function Y({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[N(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${q}`}),N(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${J}`}),N(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),N(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),N(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),N(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:X,cubicBezierEaseOut:Z}=p;function Q({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[N(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${X}`}),N(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${Z}`}),N(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),N(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),N(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),N(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:$,cubicBezierEaseOut:ee}=p;function te({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[N(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${$}`}),N(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${ee}`}),N(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),N(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),N(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),N(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var ne=N([R(`drawer`,`
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
 `,[Q(),Y(),te(),K(),E(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),E(`native-scrollbar`,[R(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),T(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[E(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),R(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),R(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[E(`native-scrollbar`,[R(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),R(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),R(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),R(`drawer-header`,`
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
 `)]),R(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),E(`right-placement`,`
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
 `)]),E(`left-placement`,`
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
 `)]),E(`top-placement`,`
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
 `)]),E(`bottom-placement`,`
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
 `)])]),N(`body`,[N(`>`,[R(`drawer-container`,`
 position: fixed;
 `)])]),R(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[N(`> *`,`
 pointer-events: all;
 `)]),R(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[E(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),k({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),re=l({name:`Drawer`,inheritAttrs:!1,props:Object.assign(Object.assign({},d.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),setup(e){let{mergedClsPrefixRef:t,namespaceRef:n,inlineThemeDisabled:r}=x(e),i=v(),a=d(`Drawer`,`-drawer`,ne,H,e,t),o=O(e.defaultWidth),s=O(e.defaultHeight),l=y(B(e,`width`),o),f=y(B(e,`height`),s),p=g(()=>{let{placement:t}=e;return t===`top`||t===`bottom`?``:A(l.value)}),m=g(()=>{let{placement:t}=e;return t===`left`||t===`right`?``:A(f.value)}),h=t=>{let{onUpdateWidth:n,"onUpdate:width":r}=e;n&&P(n,t),r&&P(r,t),o.value=t},_=t=>{let{onUpdateHeight:n,"onUpdate:width":r}=e;n&&P(n,t),r&&P(r,t),s.value=t},b=g(()=>[{width:p.value,height:m.value},e.drawerStyle||``]);function S(t){let{onMaskClick:n,maskClosable:r}=e;r&&E(!1),n&&n(t)}function C(e){S(e)}let w=j();function T(t){var n;(n=e.onEsc)==null||n.call(e),e.show&&e.closeOnEsc&&F(t)&&(w.value||E(!1))}function E(t){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=e;r&&P(r,t),i&&P(i,t),n&&!t&&P(n,t)}L(u,{isMountedRef:i,mergedThemeRef:a,mergedClsPrefixRef:t,doUpdateShow:E,doUpdateHeight:_,doUpdateWidth:h});let D=g(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:o,lineHeight:s,headerPadding:c,footerPadding:l,borderRadius:u,bodyPadding:d,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=a.value;return{"--n-line-height":s,"--n-color":r,"--n-border-radius":u,"--n-text-color":i,"--n-box-shadow":o,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":c,"--n-body-padding":d,"--n-footer-padding":l,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),k=r?c(`drawer`,void 0,D,e):void 0;return{mergedClsPrefix:t,namespace:n,mergedBodyStyle:b,handleOutsideClick:C,handleMaskClick:S,handleEsc:T,mergedTheme:a,cssVars:r?void 0:D,themeClass:k?.themeClass,onRender:k?.onRender,isMounted:i}},render(){let{mergedClsPrefix:t}=this;return I(_,{to:this.to,show:this.show},{default:()=>{var n;return(n=this.onRender)==null||n.call(this),w(I(`div`,{class:[`${t}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:`none`},this.showMask?I(e,{name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?I(`div`,{"aria-hidden":!0,class:[`${t}-drawer-mask`,this.showMask===`transparent`&&`${t}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,I(U,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[r,{zIndex:this.zIndex,enabled:this.show}]])}})}}),ie=l({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let e=s(u,null);e||z(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:t}=e;function n(){t(!1)}return{handleCloseClick:n,mergedTheme:e.mergedThemeRef,mergedClsPrefix:e.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:t,nativeScrollbar:n,mergedTheme:r,bodyClass:i,bodyStyle:a,bodyContentClass:o,bodyContentStyle:s,headerClass:c,headerStyle:l,footerClass:u,footerStyle:d,scrollbarProps:f,closable:p,$slots:m}=this;return I(`div`,{role:`none`,class:[`${t}-drawer-content`,n&&`${t}-drawer-content--native-scrollbar`]},m.header||e||p?I(`div`,{class:[`${t}-drawer-header`,c],style:l,role:`none`},I(`div`,{class:`${t}-drawer-header__main`,role:`heading`,"aria-level":`1`},m.header===void 0?e:m.header()),p&&I(S,{onClick:this.handleCloseClick,clsPrefix:t,class:`${t}-drawer-header__close`,absolute:!0})):null,n?I(`div`,{class:[`${t}-drawer-body`,i],style:a,role:`none`},I(`div`,{class:[`${t}-drawer-body-content-wrapper`,o],style:s,role:`none`},m)):I(D,Object.assign({themeOverrides:r.peerOverrides.Scrollbar,theme:r.peers.Scrollbar},f,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,o],contentStyle:s}),m),m.footer?I(`div`,{class:[`${t}-drawer-footer`,u],style:d,role:`none`},m.footer()):null)}});export{re as n,ie as t};