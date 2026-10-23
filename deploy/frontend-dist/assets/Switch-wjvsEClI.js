import{Gi as e,Jn as t,Ki as n,Kn as r,Li as i,On as a,Qn as o,Qr as s,Ur as c,_a as l,bn as u,ci as d,cr as f,di as p,fi as m,gn as h,ha as g,ii as _,jn as v,li as y,mn as b,ni as x,pi as S,qn as C,tr as w,ui as T}from"./router-pJPPzlzY.js";import{t as E}from"./index-43w-3aur.js";function D(e){let{primaryColor:t,opacityDisabled:n,borderRadius:r,textColor3:i}=e;return Object.assign(Object.assign({},E),{iconColor:i,textColor:`white`,loadingColor:t,opacityDisabled:n,railColor:`rgba(0, 0, 0, .14)`,railColorActive:t,buttonBoxShadow:`0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)`,buttonColor:`#FFF`,railBorderRadiusSmall:r,railBorderRadiusMedium:r,railBorderRadiusLarge:r,buttonBorderRadiusSmall:r,buttonBorderRadiusMedium:r,buttonBorderRadiusLarge:r,boxShadowFocus:`0 0 0 2px ${s(t,{alpha:.2})}`})}var O={name:`Switch`,common:b,self:D},k=y(`switch`,`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[T(`children-placeholder`,`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),T(`rail-placeholder`,`
 display: flex;
 flex-wrap: none;
 `),T(`button-placeholder`,`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),y(`base-loading`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[u({left:`50%`,top:`50%`,originalTransform:`translateX(-50%) translateY(-50%)`})]),T(`checked, unchecked`,`
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
 `),T(`checked`,`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),T(`unchecked`,`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),d(`&:focus`,[T(`rail`,`
 box-shadow: var(--n-box-shadow-focus);
 `)]),p(`round`,[T(`rail`,`border-radius: calc(var(--n-rail-height) / 2);`,[T(`button`,`border-radius: calc(var(--n-button-height) / 2);`)])]),m(`disabled`,[m(`icon`,[p(`rubber-band`,[p(`pressed`,[T(`rail`,[T(`button`,`max-width: var(--n-button-width-pressed);`)])]),T(`rail`,[d(`&:active`,[T(`button`,`max-width: var(--n-button-width-pressed);`)])]),p(`active`,[p(`pressed`,[T(`rail`,[T(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])]),T(`rail`,[d(`&:active`,[T(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])])])])])]),p(`active`,[T(`rail`,[T(`button`,`left: calc(100% - var(--n-button-width) - var(--n-offset))`)])]),T(`rail`,`
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
 `,[T(`button-icon`,`
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
 `,[u()]),T(`button`,`
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
 `)]),p(`active`,[T(`rail`,`background-color: var(--n-rail-color-active);`)]),p(`loading`,[T(`rail`,`
 cursor: wait;
 `)]),p(`disabled`,[T(`rail`,`
 cursor: not-allowed;
 opacity: .5;
 `)])]),A=Object.assign(Object.assign({},v.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]}),j,M=e({name:`Switch`,props:A,slots:Object,setup(e){j===void 0&&(j=typeof CSS<`u`?CSS.supports===void 0?!1:CSS.supports(`width`,`max(1px)`):!0);let{mergedClsPrefixRef:n,inlineThemeDisabled:a,mergedComponentPropsRef:o}=t(e),s=v(`Switch`,`-switch`,k,O,e,n),u=r(e,{mergedSize(t){return e.size===void 0?t?t.mergedSize.value:o?.value?.Switch?.size||`medium`:e.size}}),{mergedSizeRef:d,mergedDisabledRef:p}=u,m=g(e.defaultValue),h=c(l(e,`value`),m),y=i(()=>h.value===e.checkedValue),b=g(!1),w=g(!1),T=i(()=>{let{railStyle:t}=e;if(t)return t({focused:w.value,checked:y.value})});function E(t){let{"onUpdate:value":n,onChange:r,onUpdateValue:i}=e,{nTriggerFormInput:a,nTriggerFormChange:o}=u;n&&f(n,t),i&&f(i,t),r&&f(r,t),m.value=t,a(),o()}function D(){let{nTriggerFormFocus:e}=u;e()}function A(){let{nTriggerFormBlur:e}=u;e()}function M(){e.loading||p.value||(h.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue))}function N(){w.value=!0,D()}function P(){w.value=!1,A(),b.value=!1}function F(t){e.loading||p.value||t.key===` `&&(h.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue),b.value=!1)}function I(t){e.loading||p.value||t.key===` `&&(t.preventDefault(),b.value=!0)}let L=i(()=>{let{value:e}=d,{self:{opacityDisabled:t,railColor:n,railColorActive:r,buttonBoxShadow:i,buttonColor:a,boxShadowFocus:o,loadingColor:c,textColor:l,iconColor:u,[S(`buttonHeight`,e)]:f,[S(`buttonWidth`,e)]:p,[S(`buttonWidthPressed`,e)]:m,[S(`railHeight`,e)]:h,[S(`railWidth`,e)]:g,[S(`railBorderRadius`,e)]:v,[S(`buttonBorderRadius`,e)]:y},common:{cubicBezierEaseInOut:b}}=s.value,C,w,T;return j?(C=`calc((${h} - ${f}) / 2)`,w=`max(${h}, ${f})`,T=`max(${g}, calc(${g} + ${f} - ${h}))`):(C=_((x(h)-x(f))/2),w=_(Math.max(x(h),x(f))),T=x(h)>x(f)?g:_(x(g)+x(f)-x(h))),{"--n-bezier":b,"--n-button-border-radius":y,"--n-button-box-shadow":i,"--n-button-color":a,"--n-button-width":p,"--n-button-width-pressed":m,"--n-button-height":f,"--n-height":w,"--n-offset":C,"--n-opacity-disabled":t,"--n-rail-border-radius":v,"--n-rail-color":n,"--n-rail-color-active":r,"--n-rail-height":h,"--n-rail-width":g,"--n-width":T,"--n-box-shadow-focus":o,"--n-loading-color":c,"--n-text-color":l,"--n-icon-color":u}}),R=a?C(`switch`,i(()=>d.value[0]),L,e):void 0;return{handleClick:M,handleBlur:P,handleFocus:N,handleKeyup:F,handleKeydown:I,mergedRailStyle:T,pressed:b,mergedClsPrefix:n,mergedValue:h,checked:y,mergedDisabled:p,cssVars:a?void 0:L,themeClass:R?.themeClass,onRender:R?.onRender}},render(){let{mergedClsPrefix:e,mergedDisabled:t,checked:r,mergedRailStyle:i,onRender:s,$slots:c}=this;s?.();let{checked:l,unchecked:u,icon:d,"checked-icon":f,"unchecked-icon":p}=c,m=!(o(d)&&o(f)&&o(p));return n(`div`,{role:`switch`,"aria-checked":r,class:[`${e}-switch`,this.themeClass,m&&`${e}-switch--icon`,r&&`${e}-switch--active`,t&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},n(`div`,{class:`${e}-switch__rail`,"aria-hidden":`true`,style:i},w(l,t=>w(u,r=>t||r?n(`div`,{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},n(`div`,{class:`${e}-switch__rail-placeholder`},n(`div`,{class:`${e}-switch__button-placeholder`}),t),n(`div`,{class:`${e}-switch__rail-placeholder`},n(`div`,{class:`${e}-switch__button-placeholder`}),r)):null)),n(`div`,{class:`${e}-switch__button`},w(d,t=>w(f,r=>w(p,i=>n(a,null,{default:()=>this.loading?n(h,Object.assign({key:`loading`,clsPrefix:e,strokeWidth:20},this.spinProps)):this.checked&&(r||t)?n(`div`,{class:`${e}-switch__button-icon`,key:r?`checked-icon`:`icon`},r||t):!this.checked&&(i||t)?n(`div`,{class:`${e}-switch__button-icon`,key:i?`unchecked-icon`:`icon`},i||t):null})))),w(l,t=>t&&n(`div`,{key:`checked`,class:`${e}-switch__checked`},t)),w(u,t=>t&&n(`div`,{key:`unchecked`,class:`${e}-switch__unchecked`},t)))))}});export{M as t};