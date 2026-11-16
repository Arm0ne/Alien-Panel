import{$n as e,$r as t,Jn as n,Ki as r,Mn as i,Ri as a,Wr as o,Yn as s,_n as c,ai as l,di as u,fi as d,ga as f,hn as p,kn as m,li as h,lr as g,mi as _,nr as v,pi as y,qi as b,qn as x,ri as S,ui as C,va as w,xn as T}from"./router-DxRn2447.js";import{t as E}from"./index-DgU6dBrR.js";function D(e){let{primaryColor:n,opacityDisabled:r,borderRadius:i,textColor3:a}=e;return Object.assign(Object.assign({},E),{iconColor:a,textColor:`white`,loadingColor:n,opacityDisabled:r,railColor:`rgba(0, 0, 0, .14)`,railColorActive:n,buttonBoxShadow:`0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)`,buttonColor:`#FFF`,railBorderRadiusSmall:i,railBorderRadiusMedium:i,railBorderRadiusLarge:i,buttonBorderRadiusSmall:i,buttonBorderRadiusMedium:i,buttonBorderRadiusLarge:i,boxShadowFocus:`0 0 0 2px ${t(n,{alpha:.2})}`})}var O={name:`Switch`,common:p,self:D},k=C(`switch`,`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[u(`children-placeholder`,`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),u(`rail-placeholder`,`
 display: flex;
 flex-wrap: none;
 `),u(`button-placeholder`,`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),C(`base-loading`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[T({left:`50%`,top:`50%`,originalTransform:`translateX(-50%) translateY(-50%)`})]),u(`checked, unchecked`,`
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
 `),u(`checked`,`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),u(`unchecked`,`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),h(`&:focus`,[u(`rail`,`
 box-shadow: var(--n-box-shadow-focus);
 `)]),d(`round`,[u(`rail`,`border-radius: calc(var(--n-rail-height) / 2);`,[u(`button`,`border-radius: calc(var(--n-button-height) / 2);`)])]),y(`disabled`,[y(`icon`,[d(`rubber-band`,[d(`pressed`,[u(`rail`,[u(`button`,`max-width: var(--n-button-width-pressed);`)])]),u(`rail`,[h(`&:active`,[u(`button`,`max-width: var(--n-button-width-pressed);`)])]),d(`active`,[d(`pressed`,[u(`rail`,[u(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])]),u(`rail`,[h(`&:active`,[u(`button`,`left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));`)])])])])])]),d(`active`,[u(`rail`,[u(`button`,`left: calc(100% - var(--n-button-width) - var(--n-offset))`)])]),u(`rail`,`
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
 `,[u(`button-icon`,`
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
 `,[T()]),u(`button`,`
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
 `)]),d(`active`,[u(`rail`,`background-color: var(--n-rail-color-active);`)]),d(`loading`,[u(`rail`,`
 cursor: wait;
 `)]),d(`disabled`,[u(`rail`,`
 cursor: not-allowed;
 opacity: .5;
 `)])]),A=Object.assign(Object.assign({},i.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]}),j,M=r({name:`Switch`,props:A,slots:Object,setup(e){j===void 0&&(j=typeof CSS<`u`?CSS.supports===void 0?!1:CSS.supports(`width`,`max(1px)`):!0);let{mergedClsPrefixRef:t,inlineThemeDisabled:r,mergedComponentPropsRef:c}=s(e),u=i(`Switch`,`-switch`,k,O,e,t),d=x(e,{mergedSize(t){return e.size===void 0?t?t.mergedSize.value:c?.value?.Switch?.size||`medium`:e.size}}),{mergedSizeRef:p,mergedDisabledRef:m}=d,h=f(e.defaultValue),v=o(w(e,`value`),h),y=a(()=>v.value===e.checkedValue),b=f(!1),C=f(!1),T=a(()=>{let{railStyle:t}=e;if(t)return t({focused:C.value,checked:y.value})});function E(t){let{"onUpdate:value":n,onChange:r,onUpdateValue:i}=e,{nTriggerFormInput:a,nTriggerFormChange:o}=d;n&&g(n,t),i&&g(i,t),r&&g(r,t),h.value=t,a(),o()}function D(){let{nTriggerFormFocus:e}=d;e()}function A(){let{nTriggerFormBlur:e}=d;e()}function M(){e.loading||m.value||(v.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue))}function N(){C.value=!0,D()}function P(){C.value=!1,A(),b.value=!1}function F(t){e.loading||m.value||t.key===` `&&(v.value===e.checkedValue?E(e.uncheckedValue):E(e.checkedValue),b.value=!1)}function I(t){e.loading||m.value||t.key===` `&&(t.preventDefault(),b.value=!0)}let L=a(()=>{let{value:e}=p,{self:{opacityDisabled:t,railColor:n,railColorActive:r,buttonBoxShadow:i,buttonColor:a,boxShadowFocus:o,loadingColor:s,textColor:c,iconColor:d,[_(`buttonHeight`,e)]:f,[_(`buttonWidth`,e)]:m,[_(`buttonWidthPressed`,e)]:h,[_(`railHeight`,e)]:g,[_(`railWidth`,e)]:v,[_(`railBorderRadius`,e)]:y,[_(`buttonBorderRadius`,e)]:b},common:{cubicBezierEaseInOut:x}}=u.value,C,w,T;return j?(C=`calc((${g} - ${f}) / 2)`,w=`max(${g}, ${f})`,T=`max(${v}, calc(${v} + ${f} - ${g}))`):(C=l((S(g)-S(f))/2),w=l(Math.max(S(g),S(f))),T=S(g)>S(f)?v:l(S(v)+S(f)-S(g))),{"--n-bezier":x,"--n-button-border-radius":b,"--n-button-box-shadow":i,"--n-button-color":a,"--n-button-width":m,"--n-button-width-pressed":h,"--n-button-height":f,"--n-height":w,"--n-offset":C,"--n-opacity-disabled":t,"--n-rail-border-radius":y,"--n-rail-color":n,"--n-rail-color-active":r,"--n-rail-height":g,"--n-rail-width":v,"--n-width":T,"--n-box-shadow-focus":o,"--n-loading-color":s,"--n-text-color":c,"--n-icon-color":d}}),R=r?n(`switch`,a(()=>p.value[0]),L,e):void 0;return{handleClick:M,handleBlur:P,handleFocus:N,handleKeyup:F,handleKeydown:I,mergedRailStyle:T,pressed:b,mergedClsPrefix:t,mergedValue:v,checked:y,mergedDisabled:m,cssVars:r?void 0:L,themeClass:R?.themeClass,onRender:R?.onRender}},render(){let{mergedClsPrefix:t,mergedDisabled:n,checked:r,mergedRailStyle:i,onRender:a,$slots:o}=this;a?.();let{checked:s,unchecked:l,icon:u,"checked-icon":d,"unchecked-icon":f}=o,p=!(e(u)&&e(d)&&e(f));return b(`div`,{role:`switch`,"aria-checked":r,class:[`${t}-switch`,this.themeClass,p&&`${t}-switch--icon`,r&&`${t}-switch--active`,n&&`${t}-switch--disabled`,this.round&&`${t}-switch--round`,this.loading&&`${t}-switch--loading`,this.pressed&&`${t}-switch--pressed`,this.rubberBand&&`${t}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},b(`div`,{class:`${t}-switch__rail`,"aria-hidden":`true`,style:i},v(s,e=>v(l,n=>e||n?b(`div`,{"aria-hidden":!0,class:`${t}-switch__children-placeholder`},b(`div`,{class:`${t}-switch__rail-placeholder`},b(`div`,{class:`${t}-switch__button-placeholder`}),e),b(`div`,{class:`${t}-switch__rail-placeholder`},b(`div`,{class:`${t}-switch__button-placeholder`}),n)):null)),b(`div`,{class:`${t}-switch__button`},v(u,e=>v(d,n=>v(f,r=>b(m,null,{default:()=>this.loading?b(c,Object.assign({key:`loading`,clsPrefix:t,strokeWidth:20},this.spinProps)):this.checked&&(n||e)?b(`div`,{class:`${t}-switch__button-icon`,key:n?`checked-icon`:`icon`},n||e):!this.checked&&(r||e)?b(`div`,{class:`${t}-switch__button-icon`,key:r?`unchecked-icon`:`icon`},r||e):null})))),v(s,e=>e&&b(`div`,{key:`checked`,class:`${t}-switch__checked`},e)),v(l,e=>e&&b(`div`,{key:`unchecked`,class:`${t}-switch__unchecked`},e)))))}});export{M as t};