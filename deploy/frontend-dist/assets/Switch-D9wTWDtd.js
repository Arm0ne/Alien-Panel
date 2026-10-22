import{Bi as e,Bn as t,Hn as n,Kn as r,Kr as i,Lr as a,Qr as o,Sn as s,Tn as c,Vn as l,Xr as u,Yn as d,ai as f,cn as p,da as m,ii as h,ji as g,la as _,mn as v,ni as y,nr as b,oi as x,ri as S,si as C,un as w,zi as T}from"./router-DMnLqP4l.js";import{t as E}from"./index-c6iM4-gc.js";function D(e){let{primaryColor:t,opacityDisabled:n,borderRadius:r,textColor3:a}=e;return Object.assign(Object.assign({},E),{iconColor:a,textColor:`white`,loadingColor:t,opacityDisabled:n,railColor:`rgba(0, 0, 0, .14)`,railColorActive:t,buttonBoxShadow:`0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)`,buttonColor:`#FFF`,railBorderRadiusSmall:r,railBorderRadiusMedium:r,railBorderRadiusLarge:r,buttonBorderRadiusSmall:r,buttonBorderRadiusMedium:r,buttonBorderRadiusLarge:r,boxShadowFocus:`0 0 0 2px ${i(t,{alpha:.2})}`})}var O={name:`Switch`,common:p,self:D},k=S(`switch`,`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[h(`children-placeholder`,`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),h(`rail-placeholder`,`
 display: flex;
 flex-wrap: none;
 `),h(`button-placeholder`,`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),S(`base-loading`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[v({left:`50%`,top:`50%`,originalTransform:`translateX(-50%) translateY(-50%)`})]),h(`checked, unchecked`,`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),h(`checked`,`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),h(`unchecked`,`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),y(`&:focus`,[h(`rail`,`
 box-shadow: var(--n-box-shadow-focus);
 `)]),f(`round`,[h(`rail`,`border-radius: calc(var(--n-rail-height) / 2);`,[h(`button`,`border-radius: calc(var(--n-button-height) / 2);`)])]),x(`disabled`,[x(`icon`,[f(`rubber-band`,[f(`pressed`,[h(`rail`,[h(`button`,`max-width: var(--n-button-width-pressed);`)])]),h(`rail`,[y(`&:active`,[h(`button`,`max-width: var(--n-button-width-pressed);`)])]),f(`active`,[f(`pressed`,[h(`rail`,[h(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])]),h(`rail`,[y(`&:active`,[h(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])])])])])]),f(`active`,[h(`rail`,[h(`button`,`left: calc(100% - var(--n-button-width) - var(--n-offset))`)])]),h(`rail`,`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[h(`button-icon`,`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[v()]),h(`button`,`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),f(`active`,[h(`rail`,`background-color: var(--n-rail-color-active);`)]),f(`loading`,[h(`rail`,`
 cursor: wait;
 `)]),f(`disabled`,[h(`rail`,`
 cursor: not-allowed;
 opacity: .5;
 `)])]),A=Object.assign(Object.assign({},c.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]}),j,M=T({name:`Switch`,props:A,slots:Object,setup(e){j===void 0&&(j=typeof CSS<`u`?CSS.supports===void 0?!1:CSS.supports(`width`,`max(1px)`):!0);let{mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:s}=n(e),d=c(`Switch`,`-switch`,k,O,e,r),f=t(e,{mergedSize(t){return e.size===void 0?t?t.mergedSize.value:s?.value?.Switch?.size||`medium`:e.size}}),{mergedSizeRef:p,mergedDisabledRef:h}=f,v=_(e.defaultValue),y=a(m(e,`value`),v),x=g(()=>y.value===e.checkedValue),S=_(!1),w=_(!1),T=g(()=>{let{railStyle:t}=e;if(t)return t({focused:w.value,checked:x.value})});function E(t){let{"onUpdate:value":n,onChange:r,onUpdateValue:i}=e,{nTriggerFormInput:a,nTriggerFormChange:o}=f;n&&b(n,t),i&&b(i,t),r&&b(r,t),v.value=t,a(),o()}function D(){let{nTriggerFormFocus:e}=f;e()}function A(){let{nTriggerFormBlur:e}=f;e()}function M(){e.loading||h.value||(y.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue))}function N(){w.value=!0,D()}function P(){w.value=!1,A(),S.value=!1}function F(t){e.loading||h.value||t.key===` `&&(y.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue),S.value=!1)}function I(t){e.loading||h.value||t.key===` `&&(t.preventDefault(),S.value=!0)}let L=g(()=>{let{value:e}=p,{self:{opacityDisabled:t,railColor:n,railColorActive:r,buttonBoxShadow:i,buttonColor:a,boxShadowFocus:s,loadingColor:c,textColor:l,iconColor:f,[C(`buttonHeight`,e)]:m,[C(`buttonWidth`,e)]:h,[C(`buttonWidthPressed`,e)]:g,[C(`railHeight`,e)]:_,[C(`railWidth`,e)]:v,[C(`railBorderRadius`,e)]:y,[C(`buttonBorderRadius`,e)]:b},common:{cubicBezierEaseInOut:x}}=d.value,S,w,T;return j?(S=`calc((${_} - ${m}) / 2)`,w=`max(${_}, ${m})`,T=`max(${v}, calc(${v} + ${m} - ${_}))`):(S=o((u(_)-u(m))/2),w=o(Math.max(u(_),u(m))),T=u(_)>u(m)?v:o(u(v)+u(m)-u(_))),{"--n-bezier":x,"--n-button-border-radius":b,"--n-button-box-shadow":i,"--n-button-color":a,"--n-button-width":h,"--n-button-width-pressed":g,"--n-button-height":m,"--n-height":w,"--n-offset":S,"--n-opacity-disabled":t,"--n-rail-border-radius":y,"--n-rail-color":n,"--n-rail-color-active":r,"--n-rail-height":_,"--n-rail-width":v,"--n-width":T,"--n-box-shadow-focus":s,"--n-loading-color":c,"--n-text-color":l,"--n-icon-color":f}}),R=i?l(`switch`,g(()=>p.value[0]),L,e):void 0;return{handleClick:M,handleBlur:P,handleFocus:N,handleKeyup:F,handleKeydown:I,mergedRailStyle:T,pressed:S,mergedClsPrefix:r,mergedValue:y,checked:x,mergedDisabled:h,cssVars:i?void 0:L,themeClass:R?.themeClass,onRender:R?.onRender}},render(){let{mergedClsPrefix:t,mergedDisabled:n,checked:i,mergedRailStyle:a,onRender:o,$slots:c}=this;o?.();let{checked:l,unchecked:u,icon:f,"checked-icon":p,"unchecked-icon":m}=c,h=!(r(f)&&r(p)&&r(m));return e(`div`,{role:`switch`,"aria-checked":i,class:[`${t}-switch`,this.themeClass,h&&`${t}-switch--icon`,i&&`${t}-switch--active`,n&&`${t}-switch--disabled`,this.round&&`${t}-switch--round`,this.loading&&`${t}-switch--loading`,this.pressed&&`${t}-switch--pressed`,this.rubberBand&&`${t}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},e(`div`,{class:`${t}-switch__rail`,"aria-hidden":`true`,style:a},d(l,n=>d(u,r=>n||r?e(`div`,{"aria-hidden":!0,class:`${t}-switch__children-placeholder`},e(`div`,{class:`${t}-switch__rail-placeholder`},e(`div`,{class:`${t}-switch__button-placeholder`}),n),e(`div`,{class:`${t}-switch__rail-placeholder`},e(`div`,{class:`${t}-switch__button-placeholder`}),r)):null)),e(`div`,{class:`${t}-switch__button`},d(f,n=>d(p,r=>d(m,i=>e(s,null,{default:()=>this.loading?e(w,Object.assign({key:`loading`,clsPrefix:t,strokeWidth:20},this.spinProps)):this.checked&&(r||n)?e(`div`,{class:`${t}-switch__button-icon`,key:r?`checked-icon`:`icon`},r||n):!this.checked&&(i||n)?e(`div`,{class:`${t}-switch__button-icon`,key:i?`unchecked-icon`:`icon`},i||n):null})))),d(l,n=>n&&e(`div`,{key:`checked`,class:`${t}-switch__checked`},n)),d(u,n=>n&&e(`div`,{key:`unchecked`,class:`${t}-switch__unchecked`},n)))))}});export{M as t};