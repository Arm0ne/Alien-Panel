import{$r as e,An as t,Cr as n,Dn as r,Gi as i,Hi as a,Hn as o,Lr as s,Mn as c,Ni as l,Nr as u,Oi as d,Or as f,Pr as p,Qi as m,Qn as h,Qr as g,Rr as _,Sr as v,Ti as y,Tn as b,Ui as x,Un as ee,Ur as S,Vi as C,Wn as w,Xt as te,Yi as T,Yn as E,Zn as D,Zt as O,ai as k,an as ne,br as A,ci as j,cr as re,da as M,ea as N,fn as P,gr as F,hr as I,ia as ie,ii as L,ir as R,ji as z,li as B,ni as V,nn as H,oa as U,oi as W,on as ae,pa as G,qi as K,qn as oe,ra as q,si as J,sn as se,ti as ce,tn as le,ur as ue,vr as de,wr as Y,yr as X,zr as fe}from"./router-P--YcraN.js";import{o as pe,t as Z}from"./focus-detector-DG7JZTzp.js";import{a as me,r as he,t as ge}from"./light-CIDkKv6a.js";import{n as _e}from"./data-freshness-DmMzYK1p.js";var Q=`v-hidden`,ve=de(`[v-hidden]`,{display:`none!important`}),ye=C({name:`Overflow`,props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateCount:Function,onUpdateOverflow:Function},setup(e,{slots:t}){let n=M(null),r=M(null);function i(i){let{value:a}=n,{getCounter:o,getTail:s}=e,c;if(c=o===void 0?r.value:o(),!a||!c)return;c.hasAttribute(Q)&&c.removeAttribute(Q);let{children:l}=a;if(i.showAllItemsBeforeCalculate)for(let e of l)e.hasAttribute(Q)&&e.removeAttribute(Q);let u=a.offsetWidth,d=[],f=t.tail?s?.():null,p=f?f.offsetWidth:0,m=!1,h=a.children.length-+!!t.tail;for(let t=0;t<h-1;++t){if(t<0)continue;let n=l[t];if(m){n.hasAttribute(Q)||n.setAttribute(Q,``);continue}else n.hasAttribute(Q)&&n.removeAttribute(Q);let r=n.offsetWidth;if(p+=r,d[t]=r,p>u){let{updateCounter:n}=e;for(let r=t;r>=0;--r){let i=h-1-r;n===void 0?c.textContent=`${i}`:n(i);let a=c.offsetWidth;if(p-=d[r],p+a<=u||r===0){m=!0,t=r-1,f&&(t===-1?(f.style.maxWidth=`${u-a}px`,f.style.boxSizing=`border-box`):f.style.maxWidth=``);let{onUpdateCount:n}=e;n&&n(i);break}}}}let{onUpdateOverflow:g}=e;m?g!==void 0&&g(!0):(g!==void 0&&g(!1),c.setAttribute(Q,``))}let a=A();return ve.mount({id:`vueuc/overflow`,head:!0,anchorMetaName:X,ssr:a}),T(()=>i({showAllItemsBeforeCalculate:!1})),{selfRef:n,counterRef:r,sync:i}},render(){let{$slots:e}=this;return i(()=>this.sync({showAllItemsBeforeCalculate:!1})),a(`div`,{class:`v-overflow`,ref:`selfRef`},[N(e,`default`),e.counter?e.counter():a(`span`,{style:{display:`inline-block`},ref:`counterRef`}),e.tail?e.tail():null])}});function be(e,t){t&&(T(()=>{let{value:n}=e;n&&I.registerHandler(n,t)}),q(e,(e,t)=>{t&&I.unregisterHandler(t)},{deep:!1}),K(()=>{let{value:t}=e;t&&I.unregisterHandler(t)}))}function $(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var xe=C({name:`Checkmark`,render(){return a(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},a(`g`,{fill:`none`},a(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Se=C({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=x(p);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,o=r?.(i),s=t?t(i,!1):h(i[this.labelField],i,!1),c=a(`div`,Object.assign({},o,{class:[`${e}-base-select-group-header`,o?.class]}),s);return i.render?i.render({node:c,option:i}):n?n({node:c,option:i,selected:!1}):c}});function Ce(e,t){return a(y,{name:`fade-in-scale-up-transition`},{default:()=>e?a(b,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>a(xe)}):null})}var we=C({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:f}=x(p),m=S(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function h(t){let{tmNode:n}=e;n.disabled||d(t,n)}function g(t){let{tmNode:n}=e;n.disabled||f(t,n)}function _(t){let{tmNode:n}=e,{value:r}=m;n.disabled||r||f(t,n)}return{multiple:r,isGrouped:S(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:m,isSelected:S(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}else return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:_,handleMouseEnter:g,handleClick:h}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:o,nodeProps:s,renderOption:c,renderLabel:l,handleClick:u,handleMouseEnter:d,handleMouseMove:f}=this,p=Ce(n,e),m=l?[l(t,n),o&&p]:[h(t[this.labelField],t,n),o&&p],g=s?.(t),_=a(`div`,Object.assign({},g,{class:[`${e}-base-select-option`,t.class,g?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:o}],style:[g?.style||``,t.style||``],onClick:$([u,g?.onClick]),onMouseenter:$([d,g?.onMouseenter]),onMousemove:$([f,g?.onMousemove])}),a(`div`,{class:`${e}-base-select-option__content`},m));return t.render?t.render({node:_,option:t,selected:n}):c?c({node:_,option:t,selected:n}):_}}),Te=k(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[k(`scrollbar`,`
 max-height: var(--n-height);
 `),k(`virtual-list`,`
 max-height: var(--n-height);
 `),k(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[W(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),k(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),k(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),W(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),W(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),W(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),W(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),k(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),k(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[J(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),L(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),L(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),J(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),J(`pending`,[L(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),J(`selected`,`
 color: var(--n-option-text-color-active);
 `,[L(`&::before`,`
 background-color: var(--n-option-color-active);
 `),J(`pending`,[L(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),J(`disabled`,`
 cursor: not-allowed;
 `,[j(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),J(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),W(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[le({enterScale:`0.5`})])])]),Ee=C({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},r.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(n){let{mergedClsPrefixRef:a,mergedRtlRef:o,mergedComponentPropsRef:s}=w(n),c=t(`InternalSelectMenu`,o,a),d=r(`InternalSelectMenu`,`-internal-select-menu`,Te,me,n,G(n,`clsPrefix`)),f=M(null),h=M(null),_=M(null),v=l(()=>n.treeMate.getFlattenedNodes()),y=l(()=>ae(v.value)),b=M(null);function x(){let{treeMate:e}=n,t=null,{value:r}=n;r===null?t=e.getFirstAvailableNode():(t=n.multiple?e.getNode((r||[])[(r||[]).length-1]):e.getNode(r),(!t||t.disabled)&&(t=e.getFirstAvailableNode())),H(t||null)}function S(){let{value:e}=b;e&&!n.treeMate.getNode(e.key)&&(b.value=null)}let C;q(()=>n.show,e=>{e?C=q(()=>n.treeMate,()=>{n.resetMenuOnOptionsChange?(n.autoPending?x():S(),i(U)):S()},{immediate:!0}):C?.()},{immediate:!0}),K(()=>{C?.()});let te=l(()=>g(d.value.self[B(`optionHeight`,n.size)])),E=l(()=>e(d.value.self[B(`padding`,n.size)])),D=l(()=>n.multiple&&Array.isArray(n.value)?new Set(n.value):new Set),O=l(()=>{let e=v.value;return e&&e.length===0}),k=l(()=>s?.value?.Select?.renderEmpty);function ne(e){let{onToggle:t}=n;t&&t(e)}function A(e){let{onScroll:t}=n;t&&t(e)}function j(e){var t;(t=_.value)==null||t.sync(),A(e)}function re(){var e;(e=_.value)==null||e.sync()}function N(){let{value:e}=b;return e||null}function P(e,t){t.disabled||H(t,!1)}function F(e,t){t.disabled||ne(t)}function I(e){var t;V(e,`action`)||(t=n.onKeyup)==null||t.call(n,e)}function ie(e){var t;V(e,`action`)||(t=n.onKeydown)==null||t.call(n,e)}function L(e){var t;(t=n.onMousedown)==null||t.call(n,e),!n.focusable&&e.preventDefault()}function R(){let{value:e}=b;e&&H(e.getNext({loop:!0}),!0)}function z(){let{value:e}=b;e&&H(e.getPrev({loop:!0}),!0)}function H(e,t=!1){b.value=e,t&&U()}function U(){var e,t;let r=b.value;if(!r)return;let i=y.value(r.key);i!==null&&(n.virtualScroll?(e=h.value)==null||e.scrollTo({index:i}):(t=_.value)==null||t.scrollTo({index:i,elSize:te.value}))}function W(e){var t;f.value?.contains(e.target)&&((t=n.onFocus)==null||t.call(n,e))}function oe(e){var t;f.value?.contains(e.relatedTarget)||(t=n.onBlur)==null||t.call(n,e)}m(p,{handleOptionMouseEnter:P,handleOptionClick:F,valueSetRef:D,pendingTmNodeRef:b,nodePropsRef:G(n,`nodeProps`),showCheckmarkRef:G(n,`showCheckmark`),multipleRef:G(n,`multiple`),valueRef:G(n,`value`),renderLabelRef:G(n,`renderLabel`),renderOptionRef:G(n,`renderOption`),labelFieldRef:G(n,`labelField`),valueFieldRef:G(n,`valueField`)}),m(u,f),T(()=>{let{value:e}=_;e&&e.sync()});let J=l(()=>{let{size:t}=n,{common:{cubicBezierEaseInOut:r},self:{height:i,borderRadius:a,color:o,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:f,optionTextColorActive:p,optionOpacityDisabled:m,optionCheckColor:h,actionTextColor:g,optionColorPending:_,optionColorActive:v,loadingColor:y,loadingSize:b,optionColorActivePending:x,[B(`optionFontSize`,t)]:ee,[B(`optionHeight`,t)]:S,[B(`optionPadding`,t)]:C}}=d.value;return{"--n-height":i,"--n-action-divider-color":c,"--n-action-text-color":g,"--n-bezier":r,"--n-border-radius":a,"--n-color":o,"--n-option-font-size":ee,"--n-group-header-text-color":s,"--n-option-check-color":h,"--n-option-color-pending":_,"--n-option-color-active":v,"--n-option-color-active-pending":x,"--n-option-height":S,"--n-option-opacity-disabled":m,"--n-option-text-color":u,"--n-option-text-color-active":p,"--n-option-text-color-disabled":f,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":e(C,`left`),"--n-option-padding-right":e(C,`right`),"--n-loading-color":y,"--n-loading-size":b}}),{inlineThemeDisabled:se}=n,ce=se?ee(`internal-select-menu`,l(()=>n.size[0]),J,n):void 0,le={selfRef:f,next:R,prev:z,getPendingTmNode:N};return be(f,n.onResize),Object.assign({mergedTheme:d,mergedClsPrefix:a,rtlEnabled:c,virtualListRef:h,scrollbarRef:_,itemSize:te,padding:E,flattenedNodes:v,empty:O,mergedRenderEmpty:k,virtualListContainer(){let{value:e}=h;return e?.listElRef},virtualListContent(){let{value:e}=h;return e?.itemsElRef},doScroll:A,handleFocusin:W,handleFocusout:oe,handleKeyUp:I,handleKeyDown:ie,handleMouseDown:L,handleVirtualListResize:re,handleVirtualListScroll:j,cssVars:se?void 0:J,themeClass:ce?.themeClass,onRender:ce?.onRender},le)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:r,themeClass:i,onRender:o}=this;return o?.(),a(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,i,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},D(e.header,e=>e&&a(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?a(`div`,{class:`${n}-base-select-menu__loading`},a(P,{clsPrefix:n,strokeWidth:20})):this.empty?a(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},E(e.empty,()=>[this.mergedRenderEmpty?.call(this)||a(H,{theme:r.peers.Empty,themeOverrides:r.peerOverrides.Empty,size:this.size})])):a(se,Object.assign({ref:`scrollbarRef`,theme:r.peers.Scrollbar,themeOverrides:r.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?a(pe,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?a(Se,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:a(we,{clsPrefix:n,key:e.key,tmNode:e})}):a(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?a(Se,{key:e.key,clsPrefix:n,tmNode:e}):a(we,{clsPrefix:n,key:e.key,tmNode:e})))}),D(e.action,e=>e&&[a(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),a(Z,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),De=L([k(`base-selection`,`
 --n-padding-single: var(--n-padding-single-top) var(--n-padding-single-right) var(--n-padding-single-bottom) var(--n-padding-single-left);
 --n-padding-multiple: var(--n-padding-multiple-top) var(--n-padding-multiple-right) var(--n-padding-multiple-bottom) var(--n-padding-multiple-left);
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[k(`base-loading`,`
 color: var(--n-loading-color);
 `),k(`base-selection-tags`,`min-height: var(--n-height);`),W(`border, state-border`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),W(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),k(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[W(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),k(`base-selection-overlay`,`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[W(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),k(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[W(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),k(`base-selection-tags`,`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),k(`base-selection-label`,`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[k(`base-selection-input`,`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[W(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),W(`render-label`,`
 color: var(--n-text-color);
 `)]),j(`disabled`,[L(`&:hover`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),J(`focus`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),J(`active`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),k(`base-selection-label`,`background-color: var(--n-color-active);`),k(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),J(`disabled`,`cursor: not-allowed;`,[W(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),k(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[k(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),W(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),k(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),k(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),k(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[W(`input`,`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),W(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>J(`${e}-status`,[W(`state-border`,`border: var(--n-border-${e});`),j(`disabled`,[L(`&:hover`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),J(`active`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),k(`base-selection-label`,`background-color: var(--n-color-active-${e});`),k(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),J(`focus`,[W(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),k(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),k(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[L(`&:last-child`,`padding-right: 0;`),k(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[W(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Oe=C({name:`InternalSelection`,props:Object.assign(Object.assign({},r.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(n){let{mergedClsPrefixRef:a,mergedRtlRef:o}=w(n),s=t(`InternalSelection`,o,a),c=M(null),u=M(null),d=M(null),f=M(null),p=M(null),m=M(null),g=M(null),_=M(null),v=M(null),y=M(null),b=M(!1),x=M(!1),S=M(!1),C=r(`InternalSelection`,`-internal-selection`,De,he,n,G(n,`clsPrefix`)),te=l(()=>n.clearable&&!n.disabled&&(S.value||n.active)),E=l(()=>n.selectedOption?n.renderTag?n.renderTag({option:n.selectedOption,handleClose:()=>{}}):n.renderLabel?n.renderLabel(n.selectedOption,!0):h(n.selectedOption[n.labelField],n.selectedOption,!0):n.placeholder),D=l(()=>{let e=n.selectedOption;if(e)return e[n.labelField]}),O=l(()=>n.multiple?!!(Array.isArray(n.selectedOptions)&&n.selectedOptions.length):n.selectedOption!==null);function k(){var e;let{value:t}=c;if(t){let{value:r}=u;r&&(r.style.width=`${t.offsetWidth}px`,n.maxTagCount!==`responsive`&&((e=v.value)==null||e.sync({showAllItemsBeforeCalculate:!1})))}}function ne(){let{value:e}=y;e&&(e.style.display=`none`)}function A(){let{value:e}=y;e&&(e.style.display=`inline-block`)}q(G(n,`active`),e=>{e||ne()}),q(G(n,`pattern`),()=>{n.multiple&&i(k)});function j(e){let{onFocus:t}=n;t&&t(e)}function re(e){let{onBlur:t}=n;t&&t(e)}function N(e){let{onDeleteOption:t}=n;t&&t(e)}function P(e){let{onClear:t}=n;t&&t(e)}function F(e){let{onPatternInput:t}=n;t&&t(e)}function I(e){(!e.relatedTarget||!d.value?.contains(e.relatedTarget))&&j(e)}function L(e){d.value?.contains(e.relatedTarget)||re(e)}function R(e){P(e)}function z(){S.value=!0}function V(){S.value=!1}function H(e){!n.active||!n.filterable||e.target!==u.value&&e.preventDefault()}function U(e){N(e)}let W=M(!1);function ae(e){if(e.key===`Backspace`&&!W.value&&!n.pattern.length){let{selectedOptions:e}=n;e?.length&&U(e[e.length-1])}}let K=null;function oe(e){let{value:t}=c;t&&(t.textContent=e.target.value,k()),n.ignoreComposition&&W.value?K=e:F(e)}function J(){W.value=!0}function se(){W.value=!1,n.ignoreComposition&&F(K),K=null}function ce(e){var t;x.value=!0,(t=n.onPatternFocus)==null||t.call(n,e)}function le(e){var t;x.value=!1,(t=n.onPatternBlur)==null||t.call(n,e)}function ue(){var e,t;if(n.filterable)x.value=!1,(e=m.value)==null||e.blur(),(t=u.value)==null||t.blur();else if(n.multiple){let{value:e}=f;e?.blur()}else{let{value:e}=p;e?.blur()}}function de(){var e,t,r;n.filterable?(x.value=!1,(e=m.value)==null||e.focus()):n.multiple?(t=f.value)==null||t.focus():(r=p.value)==null||r.focus()}function Y(){let{value:e}=u;e&&(A(),e.focus())}function X(){let{value:e}=u;e&&e.blur()}function fe(e){let{value:t}=g;t&&t.setTextContent(`+${e}`)}function pe(){let{value:e}=_;return e}function Z(){return u.value}let me=null;function ge(){me!==null&&window.clearTimeout(me)}function _e(){n.active||(ge(),me=window.setTimeout(()=>{O.value&&(b.value=!0)},100))}function Q(){ge()}function ve(e){e||(ge(),b.value=!1)}q(O,e=>{e||(b.value=!1)}),T(()=>{ie(()=>{let e=m.value;e&&(n.disabled?e.removeAttribute(`tabindex`):e.tabIndex=x.value?-1:0)})}),be(d,n.onResize);let{inlineThemeDisabled:ye}=n,$=l(()=>{let{size:t}=n,{common:{cubicBezierEaseInOut:r},self:{fontWeight:i,borderRadius:a,color:o,placeholderColor:s,textColor:c,paddingSingle:l,paddingMultiple:u,caretColor:d,colorDisabled:f,textColorDisabled:p,placeholderColorDisabled:m,colorActive:h,boxShadowFocus:g,boxShadowActive:_,boxShadowHover:v,border:y,borderFocus:b,borderHover:x,borderActive:ee,arrowColor:S,arrowColorDisabled:w,loadingColor:te,colorActiveWarning:T,boxShadowFocusWarning:E,boxShadowActiveWarning:D,boxShadowHoverWarning:O,borderWarning:k,borderFocusWarning:ne,borderHoverWarning:A,borderActiveWarning:j,colorActiveError:re,boxShadowFocusError:M,boxShadowActiveError:N,boxShadowHoverError:P,borderError:F,borderFocusError:I,borderHoverError:ie,borderActiveError:L,clearColor:R,clearColorHover:z,clearColorPressed:V,clearSize:H,arrowSize:U,[B(`height`,t)]:W,[B(`fontSize`,t)]:ae}}=C.value,G=e(l),K=e(u);return{"--n-bezier":r,"--n-border":y,"--n-border-active":ee,"--n-border-focus":b,"--n-border-hover":x,"--n-border-radius":a,"--n-box-shadow-active":_,"--n-box-shadow-focus":g,"--n-box-shadow-hover":v,"--n-caret-color":d,"--n-color":o,"--n-color-active":h,"--n-color-disabled":f,"--n-font-size":ae,"--n-height":W,"--n-padding-single-top":G.top,"--n-padding-multiple-top":K.top,"--n-padding-single-right":G.right,"--n-padding-multiple-right":K.right,"--n-padding-single-left":G.left,"--n-padding-multiple-left":K.left,"--n-padding-single-bottom":G.bottom,"--n-padding-multiple-bottom":K.bottom,"--n-placeholder-color":s,"--n-placeholder-color-disabled":m,"--n-text-color":c,"--n-text-color-disabled":p,"--n-arrow-color":S,"--n-arrow-color-disabled":w,"--n-loading-color":te,"--n-color-active-warning":T,"--n-box-shadow-focus-warning":E,"--n-box-shadow-active-warning":D,"--n-box-shadow-hover-warning":O,"--n-border-warning":k,"--n-border-focus-warning":ne,"--n-border-hover-warning":A,"--n-border-active-warning":j,"--n-color-active-error":re,"--n-box-shadow-focus-error":M,"--n-box-shadow-active-error":N,"--n-box-shadow-hover-error":P,"--n-border-error":F,"--n-border-focus-error":I,"--n-border-hover-error":ie,"--n-border-active-error":L,"--n-clear-size":H,"--n-clear-color":R,"--n-clear-color-hover":z,"--n-clear-color-pressed":V,"--n-arrow-size":U,"--n-font-weight":i}}),xe=ye?ee(`internal-selection`,l(()=>n.size[0]),$,n):void 0;return{mergedTheme:C,mergedClearable:te,mergedClsPrefix:a,rtlEnabled:s,patternInputFocused:x,filterablePlaceholder:E,label:D,selected:O,showTagsPanel:b,isComposing:W,counterRef:g,counterWrapperRef:_,patternInputMirrorRef:c,patternInputRef:u,selfRef:d,multipleElRef:f,singleElRef:p,patternInputWrapperRef:m,overflowRef:v,inputTagElRef:y,handleMouseDown:H,handleFocusin:I,handleClear:R,handleMouseEnter:z,handleMouseLeave:V,handleDeleteOption:U,handlePatternKeyDown:ae,handlePatternInputInput:oe,handlePatternInputBlur:le,handlePatternInputFocus:ce,handleMouseEnterCounter:_e,handleMouseLeaveCounter:Q,handleFocusout:L,handleCompositionEnd:se,handleCompositionStart:J,onPopoverUpdateShow:ve,focus:de,focusInput:Y,blur:ue,blurInput:X,updateCounter:fe,getCounter:pe,getTail:Z,renderLabel:n.renderLabel,cssVars:ye?void 0:$,themeClass:xe?.themeClass,onRender:xe?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:i,maxTagCount:o,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:u,renderTag:d,renderLabel:f}=this;u?.();let p=o===`responsive`,m=typeof o==`number`,g=p||m,_=a(oe,null,{default:()=>a(te,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),v;if(t){let{labelField:e}=this,t=t=>a(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},d?d({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):a(_e,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>f?f(t,!0):h(t[e],t,!0)})),s=()=>(m?this.selectedOptions.slice(0,o):this.selectedOptions).map(t),u=i?a(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},a(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),a(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,y=p?()=>a(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},a(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,b;if(m){let e=this.selectedOptions.length-o;e>0&&(b=a(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},a(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let x=p?i?a(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y,tail:()=>u}):a(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:y}):m&&b?s().concat(b):s(),ee=g?()=>a(`div`,{class:`${c}-base-selection-popover`},p?s():this.selectedOptions.map(t)):void 0,S=g?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,C=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?a(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},a(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,w=i?a(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},x,p?null:u,_):a(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},x,_);v=a(z,null,g?a(O,Object.assign({},S,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>w,default:ee}):w,C)}else if(i){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;v=a(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:re(this.label)},a(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?a(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},a(`div`,{class:`${c}-base-selection-overlay__wrapper`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):h(this.label,this.selectedOption,!0))):null,t?a(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},a(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,_)}else v=a(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?a(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},a(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):a(`div`,{class:`${c}-base-selection-input`,title:re(this.label),key:`input`},a(`div`,{class:`${c}-base-selection-input__content`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):h(this.label,this.selectedOption,!0))),_);return a(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},v,s?a(`div`,{class:`${c}-base-selection__border`}):null,s?a(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ke(e){return e.type===`group`}function Ae(e){return e.type===`ignored`}function je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Me(e,t){return{getIsGroup:ke,getIgnored:Ae,getKey(t){return ke(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function Ne(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ke(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(Ae(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function Pe(e,t,n){let r=new Map;return e.forEach(e=>{ke(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var Fe=L([k(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),k(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[le({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Ie=C({name:`Select`,props:Object.assign(Object.assign({},r.props),{to:f.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:u}=w(e),d=r(`Select`,`-select`,Fe,ge,e,t),p=M(e.defaultValue),m=fe(G(e,`value`),p),h=M(!1),g=M(``),v=s(e,[`items`,`options`]),y=M([]),b=M([]),x=l(()=>b.value.concat(y.value).concat(v.value)),S=l(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return je(e,i);let a=t[r];return typeof a==`string`?je(e,a):typeof a==`number`?je(e,String(a)):!1}}),C=l(()=>{if(e.remote)return v.value;{let{value:t}=x,{value:n}=g;return!n.length||!e.filterable?t:Ne(t,S.value,n,e.childrenField)}}),te=l(()=>{let{valueField:t,childrenField:n}=e,r=Me(t,n);return ne(C.value,r)}),T=l(()=>Pe(x.value,e.valueField,e.childrenField)),E=M(!1),D=fe(G(e,`show`),E),O=M(null),k=M(null),A=M(null),{localeRef:j}=c(`Select`),re=l(()=>e.placeholder??j.value.placeholder),N=[],P=M(new Map),F=l(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function I(t){let n=e.remote,{value:r}=P,{value:i}=T,{value:a}=F,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let ie=l(()=>{if(e.multiple){let{value:e}=m;return Array.isArray(e)?I(e):[]}return null}),L=l(()=>{let{value:t}=m;return!e.multiple&&!Array.isArray(t)?t===null?null:I([t])[0]||null:null}),z=o(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:u?.value?.Select?.size||`medium`}}),{mergedSizeRef:B,mergedDisabledRef:H,mergedStatusRef:U}=z;function W(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=z;r&&R(r,t,n),a&&R(a,t,n),i&&R(i,t,n),p.value=t,o(),s()}function ae(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=z;n&&R(n,t),r()}function K(){let{onClear:t}=e;t&&R(t)}function oe(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=z;n&&R(n,t),i(),r&&Y()}function J(t){let{onSearch:n}=e;n&&R(n,t)}function se(t){let{onScroll:n}=e;n&&R(n,t)}function le(){var t;let{remote:n,multiple:r}=e;if(n){let{value:n}=P;if(r){let{valueField:r}=e;(t=ie.value)==null||t.forEach(e=>{n.set(e[r],e)})}else{let t=L.value;t&&n.set(t[e.valueField],t)}}}function de(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&R(n,t),r&&R(r,t),E.value=t}function Y(){H.value||(de(!0),E.value=!0,e.filterable&&Le())}function X(){de(!1)}function pe(){g.value=``,b.value=N}let Z=M(!1);function me(){e.filterable&&(Z.value=!0)}function he(){e.filterable&&(Z.value=!1,D.value||pe())}function _e(){H.value||(D.value?e.filterable?Le():X():Y())}function Q(e){(A.value?.selfRef)?.contains(e.relatedTarget)||(h.value=!1,ae(e),X())}function ve(e){oe(e),h.value=!0}function ye(){h.value=!0}function be(e){O.value?.$el.contains(e.relatedTarget)||(h.value=!1,ae(e),X())}function $(){var e;(e=O.value)==null||e.focus(),X()}function xe(e){D.value&&(O.value?.$el.contains(ce(e))||X())}function Se(t){if(!Array.isArray(t))return[];if(F.value)return Array.from(t);{let{remote:n}=e,{value:r}=T;if(n){let{value:e}=P;return t.filter(t=>r.has(t)||e.has(t))}else return t.filter(e=>r.has(e))}}function Ce(e){we(e.rawNode)}function we(t){if(H.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=b,t=e[0]||null;if(t){let e=y.value;e.length?e.push(t):y.value=[t],b.value=N}}if(r&&P.value.set(t[a],t),e.multiple){let e=Se(m.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Te(t[a]);~e&&(y.value.splice(e,1),i&&(g.value=``))}}else e.push(t[a]),i&&(g.value=``);W(e,I(e))}else{if(n&&!r){let e=Te(t[a]);~e?y.value=[y.value[e]]:y.value=N}Ie(),X(),W(t[a],t)}}function Te(t){return y.value.findIndex(n=>n[e.valueField]===t)}function Ee(t){D.value||Y();let{value:n}=t.target;g.value=n;let{tag:r,remote:i}=e;if(J(n),r&&!i){if(!n){b.value=N;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;v.value.some(e=>e[i]===r[i]||e[a]===r[a])||y.value.some(e=>e[i]===r[i]||e[a]===r[a])?b.value=N:b.value=[r]}}function De(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&X(),r&&!i&&a&&(y.value=N),K(),n?W([],[]):W(null,null)}function Oe(e){!V(e,`action`)&&!V(e,`empty`)&&!V(e,`header`)&&e.preventDefault()}function ke(e){se(e)}function Ae(t){var n,r,i;if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!O.value?.isComposing){if(D.value){let t=A.value?.getPendingTmNode();t?Ce(t):e.filterable||(X(),Ie())}else if(Y(),e.tag&&Z.value){let t=b.value[0];if(t){let n=t[e.valueField],{value:r}=m;e.multiple&&Array.isArray(r)&&r.includes(n)||we(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;D.value&&((n=A.value)==null||n.prev());break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;D.value?(r=A.value)==null||r.next():Y();break;case`Escape`:D.value&&(ue(t),X()),(i=O.value)==null||i.focus();break}}function Ie(){var e;(e=O.value)==null||e.focus()}function Le(){var e;(e=O.value)==null||e.focusInput()}function Re(){var e;D.value&&((e=k.value)==null||e.syncPosition())}le(),q(G(e,`options`),le);let ze={focus:()=>{var e;(e=O.value)==null||e.focus()},focusInput:()=>{var e;(e=O.value)==null||e.focusInput()},blur:()=>{var e;(e=O.value)==null||e.blur()},blurInput:()=>{var e;(e=O.value)==null||e.blurInput()}},Be=l(()=>{let{self:{menuBoxShadow:e}}=d.value;return{"--n-menu-box-shadow":e}}),Ve=a?ee(`select`,void 0,Be,e):void 0;return Object.assign(Object.assign({},ze),{mergedStatus:U,mergedClsPrefix:t,mergedBordered:n,namespace:i,treeMate:te,isMounted:_(),triggerRef:O,menuRef:A,pattern:g,uncontrolledShow:E,mergedShow:D,adjustedTo:f(e),uncontrolledValue:p,mergedValue:m,followerRef:k,localizedPlaceholder:re,selectedOption:L,selectedOptions:ie,mergedSize:B,mergedDisabled:H,focused:h,activeWithoutMenuOpen:Z,inlineThemeDisabled:a,onTriggerInputFocus:me,onTriggerInputBlur:he,handleTriggerOrMenuResize:Re,handleMenuFocus:ye,handleMenuBlur:be,handleMenuTabOut:$,handleTriggerClick:_e,handleToggle:Ce,handleDeleteOption:we,handlePatternInput:Ee,handleClear:De,handleTriggerBlur:Q,handleTriggerFocus:ve,handleKeydown:Ae,handleMenuAfterLeave:pe,handleMenuClickOutside:xe,handleMenuScroll:ke,handleMenuKeydown:Ae,handleMenuMousedown:Oe,mergedTheme:d,cssVars:a?void 0:Be,themeClass:Ve?.themeClass,onRender:Ve?.onRender})},render(){return a(`div`,{class:`${this.mergedClsPrefix}-select`},a(Y,null,{default:()=>[a(n,null,{default:()=>a(Oe,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),a(F,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===f.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>a(y,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),U(a(Ee,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[d,this.mergedShow],[v,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[v,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{$ as i,Me as n,Ee as r,Ie as t};