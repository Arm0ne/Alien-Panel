import{Bn as e,Cn as t,Fr as n,Jr as r,Li as i,Ri as a,Rn as o,Wn as s,Wr as c,Xr as l,ai as u,bn as d,cn as f,ei as p,er as m,fn as h,ii as g,ki as _,la as v,ni as y,on as b,qn as x,ri as S,sa as C,ti as w,zn as T}from"./router-BEJlwQm6.js";import{t as E}from"./index-DK8sHnu3.js";function D(e){let{primaryColor:t,opacityDisabled:n,borderRadius:r,textColor3:i}=e;return Object.assign(Object.assign({},E),{iconColor:i,textColor:`white`,loadingColor:t,opacityDisabled:n,railColor:`rgba(0, 0, 0, .14)`,railColorActive:t,buttonBoxShadow:`0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)`,buttonColor:`#FFF`,railBorderRadiusSmall:r,railBorderRadiusMedium:r,railBorderRadiusLarge:r,buttonBorderRadiusSmall:r,buttonBorderRadiusMedium:r,buttonBorderRadiusLarge:r,boxShadowFocus:`0 0 0 2px ${c(t,{alpha:.2})}`})}var O={name:`Switch`,common:b,self:D},k=w(`switch`,`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[y(`children-placeholder`,`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),y(`rail-placeholder`,`
 display: flex;
 flex-wrap: none;
 `),y(`button-placeholder`,`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),w(`base-loading`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[h({left:`50%`,top:`50%`,originalTransform:`translateX(-50%) translateY(-50%)`})]),y(`checked, unchecked`,`
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
 `),y(`checked`,`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),y(`unchecked`,`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),p(`&:focus`,[y(`rail`,`
 box-shadow: var(--n-box-shadow-focus);
 `)]),S(`round`,[y(`rail`,`border-radius: calc(var(--n-rail-height) / 2);`,[y(`button`,`border-radius: calc(var(--n-button-height) / 2);`)])]),g(`disabled`,[g(`icon`,[S(`rubber-band`,[S(`pressed`,[y(`rail`,[y(`button`,`max-width: var(--n-button-width-pressed);`)])]),y(`rail`,[p(`&:active`,[y(`button`,`max-width: var(--n-button-width-pressed);`)])]),S(`active`,[S(`pressed`,[y(`rail`,[y(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])]),y(`rail`,[p(`&:active`,[y(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])])])])])]),S(`active`,[y(`rail`,[y(`button`,`left: calc(100% - var(--n-button-width) - var(--n-offset))`)])]),y(`rail`,`
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
 `,[y(`button-icon`,`
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
 `,[h()]),y(`button`,`
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
 `)]),S(`active`,[y(`rail`,`background-color: var(--n-rail-color-active);`)]),S(`loading`,[y(`rail`,`
 cursor: wait;
 `)]),S(`disabled`,[y(`rail`,`
 cursor: not-allowed;
 opacity: .5;
 `)])]),A=Object.assign(Object.assign({},t.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]}),j,M=i({name:`Switch`,props:A,slots:Object,setup(i){j===void 0&&(j=typeof CSS<`u`?CSS.supports===void 0?!1:CSS.supports(`width`,`max(1px)`):!0);let{mergedClsPrefixRef:a,inlineThemeDisabled:s,mergedComponentPropsRef:c}=e(i),d=t(`Switch`,`-switch`,k,O,i,a),f=o(i,{mergedSize(e){return i.size===void 0?e?e.mergedSize.value:c?.value?.Switch?.size||`medium`:i.size}}),{mergedSizeRef:p,mergedDisabledRef:h}=f,g=C(i.defaultValue),y=n(v(i,`value`),g),b=_(()=>y.value===i.checkedValue),x=C(!1),S=C(!1),w=_(()=>{let{railStyle:e}=i;if(e)return e({focused:S.value,checked:b.value})});function E(e){let{"onUpdate:value":t,onChange:n,onUpdateValue:r}=i,{nTriggerFormInput:a,nTriggerFormChange:o}=f;t&&m(t,e),r&&m(r,e),n&&m(n,e),g.value=e,a(),o()}function D(){let{nTriggerFormFocus:e}=f;e()}function A(){let{nTriggerFormBlur:e}=f;e()}function M(){i.loading||h.value||(y.value===i.checkedValue?E(i.uncheckedValue):E(i.checkedValue))}function N(){S.value=!0,D()}function P(){S.value=!1,A(),x.value=!1}function F(e){i.loading||h.value||e.key===` `&&(y.value===i.checkedValue?E(i.uncheckedValue):E(i.checkedValue),x.value=!1)}function I(e){i.loading||h.value||e.key===` `&&(e.preventDefault(),x.value=!0)}let L=_(()=>{let{value:e}=p,{self:{opacityDisabled:t,railColor:n,railColorActive:i,buttonBoxShadow:a,buttonColor:o,boxShadowFocus:s,loadingColor:c,textColor:f,iconColor:m,[u(`buttonHeight`,e)]:h,[u(`buttonWidth`,e)]:g,[u(`buttonWidthPressed`,e)]:_,[u(`railHeight`,e)]:v,[u(`railWidth`,e)]:y,[u(`railBorderRadius`,e)]:b,[u(`buttonBorderRadius`,e)]:x},common:{cubicBezierEaseInOut:S}}=d.value,C,w,T;return j?(C=`calc((${v} - ${h}) / 2)`,w=`max(${v}, ${h})`,T=`max(${y}, calc(${y} + ${h} - ${v}))`):(C=l((r(v)-r(h))/2),w=l(Math.max(r(v),r(h))),T=r(v)>r(h)?y:l(r(y)+r(h)-r(v))),{"--n-bezier":S,"--n-button-border-radius":x,"--n-button-box-shadow":a,"--n-button-color":o,"--n-button-width":g,"--n-button-width-pressed":_,"--n-button-height":h,"--n-height":w,"--n-offset":C,"--n-opacity-disabled":t,"--n-rail-border-radius":b,"--n-rail-color":n,"--n-rail-color-active":i,"--n-rail-height":v,"--n-rail-width":y,"--n-width":T,"--n-box-shadow-focus":s,"--n-loading-color":c,"--n-text-color":f,"--n-icon-color":m}}),R=s?T(`switch`,_(()=>p.value[0]),L,i):void 0;return{handleClick:M,handleBlur:P,handleFocus:N,handleKeyup:F,handleKeydown:I,mergedRailStyle:w,pressed:x,mergedClsPrefix:a,mergedValue:y,checked:b,mergedDisabled:h,cssVars:s?void 0:L,themeClass:R?.themeClass,onRender:R?.onRender}},render(){let{mergedClsPrefix:e,mergedDisabled:t,checked:n,mergedRailStyle:r,onRender:i,$slots:o}=this;i?.();let{checked:c,unchecked:l,icon:u,"checked-icon":p,"unchecked-icon":m}=o,h=!(s(u)&&s(p)&&s(m));return a(`div`,{role:`switch`,"aria-checked":n,class:[`${e}-switch`,this.themeClass,h&&`${e}-switch--icon`,n&&`${e}-switch--active`,t&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},a(`div`,{class:`${e}-switch__rail`,"aria-hidden":`true`,style:r},x(c,t=>x(l,n=>t||n?a(`div`,{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},a(`div`,{class:`${e}-switch__rail-placeholder`},a(`div`,{class:`${e}-switch__button-placeholder`}),t),a(`div`,{class:`${e}-switch__rail-placeholder`},a(`div`,{class:`${e}-switch__button-placeholder`}),n)):null)),a(`div`,{class:`${e}-switch__button`},x(u,t=>x(p,n=>x(m,r=>a(d,null,{default:()=>this.loading?a(f,Object.assign({key:`loading`,clsPrefix:e,strokeWidth:20},this.spinProps)):this.checked&&(n||t)?a(`div`,{class:`${e}-switch__button-icon`,key:n?`checked-icon`:`icon`},n||t):!this.checked&&(r||t)?a(`div`,{class:`${e}-switch__button-icon`,key:r?`unchecked-icon`:`icon`},r||t):null})))),x(c,t=>t&&a(`div`,{key:`checked`,class:`${e}-switch__checked`},t)),x(l,t=>t&&a(`div`,{key:`unchecked`,class:`${e}-switch__unchecked`},t)))))}});export{M as t};