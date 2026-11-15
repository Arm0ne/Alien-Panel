import{Ai as e,An as t,Cr as n,Dr as r,Fn as i,Hr as a,Ii as o,Ji as s,Jn as c,Jr as l,Ki as u,Ln as d,Mn as f,Ni as p,Nr as m,Or as h,Qi as g,Qn as _,Ri as v,Rr as y,Tr as b,Ur as x,Wr as S,Xi as C,Yn as ee,_n as te,aa as ne,br as w,ca as T,da as E,di as D,dn as O,ea as k,er as A,fi as j,fn as re,ga as M,hr as ie,ii as N,kr as P,la as ae,li as F,lr as I,mi as L,nn as R,nr as z,oi as B,on as V,pi as H,pr as U,qi as W,qn as G,ra as K,ri as oe,rr as q,si as se,sn as ce,tn as le,ui as J,un as ue,va as Y,wr as de,xr as X,zr as Z}from"./router-XWKC2qFx.js";import{o as fe,t as pe}from"./focus-detector-DX4Oths1.js";import{a as me,r as he,t as ge}from"./light-TRE-snx9.js";import{n as _e}from"./data-freshness-D9miZ83X.js";var Q=`v-hidden`,ve=n(`[v-hidden]`,{display:`none!important`}),ye=u({name:`Overflow`,props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateCount:Function,onUpdateOverflow:Function},setup(e,{slots:t}){let n=M(null),r=M(null);function i(i){let{value:a}=n,{getCounter:o,getTail:s}=e,c;if(c=o===void 0?r.value:o(),!a||!c)return;c.hasAttribute(Q)&&c.removeAttribute(Q);let{children:l}=a;if(i.showAllItemsBeforeCalculate)for(let e of l)e.hasAttribute(Q)&&e.removeAttribute(Q);let u=a.offsetWidth,d=[],f=t.tail?s?.():null,p=f?f.offsetWidth:0,m=!1,h=a.children.length-+!!t.tail;for(let t=0;t<h-1;++t){if(t<0)continue;let n=l[t];if(m){n.hasAttribute(Q)||n.setAttribute(Q,``);continue}else n.hasAttribute(Q)&&n.removeAttribute(Q);let r=n.offsetWidth;if(p+=r,d[t]=r,p>u){let{updateCounter:n}=e;for(let r=t;r>=0;--r){let i=h-1-r;n===void 0?c.textContent=`${i}`:n(i);let a=c.offsetWidth;if(p-=d[r],p+a<=u||r===0){m=!0,t=r-1,f&&(t===-1?(f.style.maxWidth=`${u-a}px`,f.style.boxSizing=`border-box`):f.style.maxWidth=``);let{onUpdateCount:n}=e;n&&n(i);break}}}}let{onUpdateOverflow:g}=e;m?g!==void 0&&g(!0):(g!==void 0&&g(!1),c.setAttribute(Q,``))}let a=b();return ve.mount({id:`vueuc/overflow`,head:!0,anchorMetaName:de,ssr:a}),k(()=>i({showAllItemsBeforeCalculate:!1})),{selfRef:n,counterRef:r,sync:i}},render(){let{$slots:e}=this;return C(()=>this.sync({showAllItemsBeforeCalculate:!1})),W(`div`,{class:`v-overflow`,ref:`selfRef`},[ne(e,`default`),e.counter?e.counter():W(`span`,{style:{display:`inline-block`},ref:`counterRef`}),e.tail?e.tail():null])}});function be(e,t){t&&(k(()=>{let{value:n}=e;n&&w.registerHandler(n,t)}),T(e,(e,t)=>{t&&w.unregisterHandler(t)},{deep:!1}),g(()=>{let{value:t}=e;t&&w.unregisterHandler(t)}))}function $(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var xe=u({name:`Checkmark`,render(){return W(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},W(`g`,{fill:`none`},W(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Se=u({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=s(Z);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):q(i[this.labelField],i,!1),s=W(`div`,Object.assign({},a,{class:[`${e}-base-select-group-header`,a?.class]}),o);return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function Ce(n,r){return W(e,{name:`fade-in-scale-up-transition`},{default:()=>n?W(t,{clsPrefix:r,class:`${r}-base-select-option__check`},{default:()=>W(xe)}):null})}var we=u({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:c,valueFieldRef:u,showCheckmarkRef:d,nodePropsRef:f,handleOptionClick:p,handleOptionMouseEnter:m}=s(Z),h=l(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function g(t){let{tmNode:n}=e;n.disabled||p(t,n)}function _(t){let{tmNode:n}=e;n.disabled||m(t,n)}function v(t){let{tmNode:n}=e,{value:r}=h;n.disabled||r||m(t,n)}return{multiple:r,isGrouped:l(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:d,nodeProps:f,isPending:h,isSelected:l(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[u.value];if(a){let{value:e}=i;return e.has(o)}else return n===o}),labelField:c,renderLabel:a,renderOption:o,handleMouseMove:v,handleMouseEnter:_,handleClick:g}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:r,isGrouped:i,showCheckmark:a,nodeProps:o,renderOption:s,renderLabel:c,handleClick:l,handleMouseEnter:u,handleMouseMove:d}=this,f=Ce(n,e),p=c?[c(t,n),a&&f]:[q(t[this.labelField],t,n),a&&f],m=o?.(t),h=W(`div`,Object.assign({},m,{class:[`${e}-base-select-option`,t.class,m?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:i,[`${e}-base-select-option--pending`]:r,[`${e}-base-select-option--show-checkmark`]:a}],style:[m?.style||``,t.style||``],onClick:$([l,m?.onClick]),onMouseenter:$([u,m?.onMouseenter]),onMousemove:$([d,m?.onMousemove])}),W(`div`,{class:`${e}-base-select-option__content`},p));return t.render?t.render({node:h,option:t,selected:n}):s?s({node:h,option:t,selected:n}):h}}),Te=J(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[J(`scrollbar`,`
 max-height: var(--n-height);
 `),J(`virtual-list`,`
 max-height: var(--n-height);
 `),J(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[D(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),J(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),J(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),D(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),D(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),D(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),D(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),J(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),J(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[j(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),F(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),F(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),j(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),j(`pending`,[F(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),j(`selected`,`
 color: var(--n-option-text-color-active);
 `,[F(`&::before`,`
 background-color: var(--n-option-color-active);
 `),j(`pending`,[F(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),j(`disabled`,`
 cursor: not-allowed;
 `,[H(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),j(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),D(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[V({enterScale:`0.5`})])])]),Ee=u({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},f.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n,mergedComponentPropsRef:r}=ee(e),a=i(`InternalSelectMenu`,n,t),o=f(`InternalSelectMenu`,`-internal-select-menu`,Te,me,e,Y(e,`clsPrefix`)),s=M(null),l=M(null),u=M(null),d=v(()=>e.treeMate.getFlattenedNodes()),p=v(()=>O(d.value)),m=M(null);function h(){let{treeMate:t}=e,n=null,{value:r}=e;r===null?n=t.getFirstAvailableNode():(n=e.multiple?t.getNode((r||[])[(r||[]).length-1]):t.getNode(r),(!n||n.disabled)&&(n=t.getFirstAvailableNode())),B(n||null)}function _(){let{value:t}=m;t&&!e.treeMate.getNode(t.key)&&(m.value=null)}let b;T(()=>e.show,t=>{t?b=T(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?h():_(),C(V)):_()},{immediate:!0}):b?.()},{immediate:!0}),g(()=>{b?.()});let x=v(()=>oe(o.value.self[L(`optionHeight`,e.size)])),S=v(()=>N(o.value.self[L(`padding`,e.size)])),te=v(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),ne=v(()=>{let e=d.value;return e&&e.length===0}),w=v(()=>r?.value?.Select?.renderEmpty);function E(t){let{onToggle:n}=e;n&&n(t)}function D(t){let{onScroll:n}=e;n&&n(t)}function A(e){var t;(t=u.value)==null||t.sync(),D(e)}function j(){var e;(e=u.value)==null||e.sync()}function re(){let{value:e}=m;return e||null}function ie(e,t){t.disabled||B(t,!1)}function P(e,t){t.disabled||E(t)}function ae(t){var n;se(t,`action`)||(n=e.onKeyup)==null||n.call(e,t)}function F(t){var n;se(t,`action`)||(n=e.onKeydown)==null||n.call(e,t)}function I(t){var n;(n=e.onMousedown)==null||n.call(e,t),!e.focusable&&t.preventDefault()}function R(){let{value:e}=m;e&&B(e.getNext({loop:!0}),!0)}function z(){let{value:e}=m;e&&B(e.getPrev({loop:!0}),!0)}function B(e,t=!1){m.value=e,t&&V()}function V(){var t,n;let r=m.value;if(!r)return;let i=p.value(r.key);i!==null&&(e.virtualScroll?(t=l.value)==null||t.scrollTo({index:i}):(n=u.value)==null||n.scrollTo({index:i,elSize:x.value}))}function H(t){var n;s.value?.contains(t.target)&&((n=e.onFocus)==null||n.call(e,t))}function U(t){var n;s.value?.contains(t.relatedTarget)||(n=e.onBlur)==null||n.call(e,t)}K(Z,{handleOptionMouseEnter:ie,handleOptionClick:P,valueSetRef:te,pendingTmNodeRef:m,nodePropsRef:Y(e,`nodeProps`),showCheckmarkRef:Y(e,`showCheckmark`),multipleRef:Y(e,`multiple`),valueRef:Y(e,`value`),renderLabelRef:Y(e,`renderLabel`),renderOptionRef:Y(e,`renderOption`),labelFieldRef:Y(e,`labelField`),valueFieldRef:Y(e,`valueField`)}),K(y,s),k(()=>{let{value:e}=u;e&&e.sync()});let W=v(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:a,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[L(`optionFontSize`,t)]:x,[L(`optionHeight`,t)]:S,[L(`optionPadding`,t)]:C}}=o.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":a,"--n-option-font-size":x,"--n-group-header-text-color":s,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":N(C,`left`),"--n-option-padding-right":N(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:G}=e,q=G?c(`internal-select-menu`,v(()=>e.size[0]),W,e):void 0,ce={selfRef:s,next:R,prev:z,getPendingTmNode:re};return be(s,e.onResize),Object.assign({mergedTheme:o,mergedClsPrefix:t,rtlEnabled:a,virtualListRef:l,scrollbarRef:u,itemSize:x,padding:S,flattenedNodes:d,empty:ne,mergedRenderEmpty:w,virtualListContainer(){let{value:e}=l;return e?.listElRef},virtualListContent(){let{value:e}=l;return e?.itemsElRef},doScroll:D,handleFocusin:H,handleFocusout:U,handleKeyUp:ae,handleKeyDown:F,handleMouseDown:I,handleVirtualListResize:j,handleVirtualListScroll:A,cssVars:G?void 0:W,themeClass:q?.themeClass,onRender:q?.onRender},ce)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:r,themeClass:i,onRender:a}=this;return a?.(),W(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,i,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},z(e.header,e=>e&&W(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?W(`div`,{class:`${n}-base-select-menu__loading`},W(te,{clsPrefix:n,strokeWidth:20})):this.empty?W(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},A(e.empty,()=>[this.mergedRenderEmpty?.call(this)||W(ce,{theme:r.peers.Empty,themeOverrides:r.peerOverrides.Empty,size:this.size})])):W(re,Object.assign({ref:`scrollbarRef`,theme:r.peers.Scrollbar,themeOverrides:r.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?W(fe,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?W(Se,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:W(we,{clsPrefix:n,key:e.key,tmNode:e})}):W(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?W(Se,{key:e.key,clsPrefix:n,tmNode:e}):W(we,{clsPrefix:n,key:e.key,tmNode:e})))}),z(e.action,e=>e&&[W(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),W(pe,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),De=F([J(`base-selection`,`
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
 `,[J(`base-loading`,`
 color: var(--n-loading-color);
 `),J(`base-selection-tags`,`min-height: var(--n-height);`),D(`border, state-border`,`
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
 `),D(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),J(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[D(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),J(`base-selection-overlay`,`
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
 `,[D(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),J(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[D(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),J(`base-selection-tags`,`
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
 `),J(`base-selection-label`,`
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
 `,[J(`base-selection-input`,`
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
 `,[D(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),D(`render-label`,`
 color: var(--n-text-color);
 `)]),H(`disabled`,[F(`&:hover`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),j(`focus`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),j(`active`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),J(`base-selection-label`,`background-color: var(--n-color-active);`),J(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),j(`disabled`,`cursor: not-allowed;`,[D(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),J(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[J(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),D(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),J(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),J(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),J(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[D(`input`,`
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
 `),D(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>j(`${e}-status`,[D(`state-border`,`border: var(--n-border-${e});`),H(`disabled`,[F(`&:hover`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),j(`active`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),J(`base-selection-label`,`background-color: var(--n-color-active-${e});`),J(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),j(`focus`,[D(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),J(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),J(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[F(`&:last-child`,`padding-right: 0;`),J(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[D(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Oe=u({name:`InternalSelection`,props:Object.assign(Object.assign({},f.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n}=ee(e),r=i(`InternalSelection`,n,t),a=M(null),o=M(null),s=M(null),l=M(null),u=M(null),d=M(null),p=M(null),m=M(null),h=M(null),g=M(null),_=M(!1),y=M(!1),b=M(!1),x=f(`InternalSelection`,`-internal-selection`,De,he,e,Y(e,`clsPrefix`)),S=v(()=>e.clearable&&!e.disabled&&(b.value||e.active)),te=v(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):q(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),ne=v(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),w=v(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function E(){var t;let{value:n}=a;if(n){let{value:r}=o;r&&(r.style.width=`${n.offsetWidth}px`,e.maxTagCount!==`responsive`&&((t=h.value)==null||t.sync({showAllItemsBeforeCalculate:!1})))}}function D(){let{value:e}=g;e&&(e.style.display=`none`)}function O(){let{value:e}=g;e&&(e.style.display=`inline-block`)}T(Y(e,`active`),e=>{e||D()}),T(Y(e,`pattern`),()=>{e.multiple&&C(E)});function A(t){let{onFocus:n}=e;n&&n(t)}function j(t){let{onBlur:n}=e;n&&n(t)}function re(t){let{onDeleteOption:n}=e;n&&n(t)}function ie(t){let{onClear:n}=e;n&&n(t)}function P(t){let{onPatternInput:n}=e;n&&n(t)}function F(e){(!e.relatedTarget||!s.value?.contains(e.relatedTarget))&&A(e)}function I(e){s.value?.contains(e.relatedTarget)||j(e)}function R(e){ie(e)}function z(){b.value=!0}function B(){b.value=!1}function V(t){!e.active||!e.filterable||t.target!==o.value&&t.preventDefault()}function H(e){re(e)}let U=M(!1);function W(t){if(t.key===`Backspace`&&!U.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&H(t[t.length-1])}}let G=null;function K(t){let{value:n}=a;n&&(n.textContent=t.target.value,E()),e.ignoreComposition&&U.value?G=t:P(t)}function oe(){U.value=!0}function se(){U.value=!1,e.ignoreComposition&&P(G),G=null}function ce(t){var n;y.value=!0,(n=e.onPatternFocus)==null||n.call(e,t)}function le(t){var n;y.value=!1,(n=e.onPatternBlur)==null||n.call(e,t)}function J(){var t,n;if(e.filterable)y.value=!1,(t=d.value)==null||t.blur(),(n=o.value)==null||n.blur();else if(e.multiple){let{value:e}=l;e?.blur()}else{let{value:e}=u;e?.blur()}}function ue(){var t,n,r;e.filterable?(y.value=!1,(t=d.value)==null||t.focus()):e.multiple?(n=l.value)==null||n.focus():(r=u.value)==null||r.focus()}function de(){let{value:e}=o;e&&(O(),e.focus())}function X(){let{value:e}=o;e&&e.blur()}function Z(e){let{value:t}=p;t&&t.setTextContent(`+${e}`)}function fe(){let{value:e}=m;return e}function pe(){return o.value}let me=null;function ge(){me!==null&&window.clearTimeout(me)}function _e(){e.active||(ge(),me=window.setTimeout(()=>{w.value&&(_.value=!0)},100))}function Q(){ge()}function ve(e){e||(ge(),_.value=!1)}T(w,e=>{e||(_.value=!1)}),k(()=>{ae(()=>{let t=d.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=y.value?-1:0)})}),be(s,e.onResize);let{inlineThemeDisabled:ye}=e,$=v(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:b,borderActive:S,arrowColor:C,arrowColorDisabled:ee,loadingColor:te,colorActiveWarning:ne,boxShadowFocusWarning:w,boxShadowActiveWarning:T,boxShadowHoverWarning:E,borderWarning:D,borderFocusWarning:O,borderHoverWarning:k,borderActiveWarning:A,colorActiveError:j,boxShadowFocusError:re,boxShadowActiveError:M,boxShadowHoverError:ie,borderError:P,borderFocusError:ae,borderHoverError:F,borderActiveError:I,clearColor:R,clearColorHover:z,clearColorPressed:B,clearSize:V,arrowSize:H,[L(`height`,t)]:U,[L(`fontSize`,t)]:W}}=x.value,G=N(c),K=N(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":y,"--n-border-hover":b,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":W,"--n-height":U,"--n-padding-single-top":G.top,"--n-padding-multiple-top":K.top,"--n-padding-single-right":G.right,"--n-padding-multiple-right":K.right,"--n-padding-single-left":G.left,"--n-padding-multiple-left":K.left,"--n-padding-single-bottom":G.bottom,"--n-padding-multiple-bottom":K.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":C,"--n-arrow-color-disabled":ee,"--n-loading-color":te,"--n-color-active-warning":ne,"--n-box-shadow-focus-warning":w,"--n-box-shadow-active-warning":T,"--n-box-shadow-hover-warning":E,"--n-border-warning":D,"--n-border-focus-warning":O,"--n-border-hover-warning":k,"--n-border-active-warning":A,"--n-color-active-error":j,"--n-box-shadow-focus-error":re,"--n-box-shadow-active-error":M,"--n-box-shadow-hover-error":ie,"--n-border-error":P,"--n-border-focus-error":ae,"--n-border-hover-error":F,"--n-border-active-error":I,"--n-clear-size":V,"--n-clear-color":R,"--n-clear-color-hover":z,"--n-clear-color-pressed":B,"--n-arrow-size":H,"--n-font-weight":r}}),xe=ye?c(`internal-selection`,v(()=>e.size[0]),$,e):void 0;return{mergedTheme:x,mergedClearable:S,mergedClsPrefix:t,rtlEnabled:r,patternInputFocused:y,filterablePlaceholder:te,label:ne,selected:w,showTagsPanel:_,isComposing:U,counterRef:p,counterWrapperRef:m,patternInputMirrorRef:a,patternInputRef:o,selfRef:s,multipleElRef:l,singleElRef:u,patternInputWrapperRef:d,overflowRef:h,inputTagElRef:g,handleMouseDown:V,handleFocusin:F,handleClear:R,handleMouseEnter:z,handleMouseLeave:B,handleDeleteOption:H,handlePatternKeyDown:W,handlePatternInputInput:K,handlePatternInputBlur:le,handlePatternInputFocus:ce,handleMouseEnterCounter:_e,handleMouseLeaveCounter:Q,handleFocusout:I,handleCompositionEnd:se,handleCompositionStart:oe,onPopoverUpdateShow:ve,focus:ue,focusInput:de,blur:J,blurInput:X,updateCounter:Z,getCounter:fe,getTail:pe,renderLabel:e.renderLabel,cssVars:ye?void 0:$,themeClass:xe?.themeClass,onRender:xe?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:r,filterable:i,maxTagCount:a,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:u,renderTag:d,renderLabel:f}=this;u?.();let p=a===`responsive`,m=typeof a==`number`,h=p||m,g=W(_,null,{default:()=>W(le,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),v;if(t){let{labelField:e}=this,t=t=>W(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},d?d({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):W(_e,{size:n,closable:!t.disabled,disabled:r,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>f?f(t,!0):q(t[e],t,!0)})),s=()=>(m?this.selectedOptions.slice(0,a):this.selectedOptions).map(t),u=i?W(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},W(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:r,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),W(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,_=p?()=>W(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},W(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:r})):void 0,y;if(m){let e=this.selectedOptions.length-a;e>0&&(y=W(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},W(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:r},{default:()=>`+${e}`})))}let b=p?i?W(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:_,tail:()=>u}):W(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:_}):m&&y?s().concat(y):s(),x=h?()=>W(`div`,{class:`${c}-base-selection-popover`},p?s():this.selectedOptions.map(t)):void 0,S=h?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,C=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?W(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},W(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,ee=i?W(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},b,p?null:u,g):W(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:r?void 0:0},b,g);v=W(o,null,h?W(R,Object.assign({},S,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>ee,default:x}):ee,C)}else if(i){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;v=W(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:U(this.label)},W(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:r,disabled:r,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?W(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},W(`div`,{class:`${c}-base-selection-overlay__wrapper`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):q(this.label,this.selectedOption,!0))):null,t?W(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},W(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,g)}else v=W(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?W(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},W(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):W(`div`,{class:`${c}-base-selection-input`,title:U(this.label),key:`input`},W(`div`,{class:`${c}-base-selection-input__content`},d?d({option:this.selectedOption,handleClose:()=>{}}):f?f(this.selectedOption,!0):q(this.label,this.selectedOption,!0))),g);return W(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},v,s?W(`div`,{class:`${c}-base-selection__border`}):null,s?W(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ke(e){return e.type===`group`}function Ae(e){return e.type===`ignored`}function je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Me(e,t){return{getIsGroup:ke,getIgnored:Ae,getKey(t){return ke(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function Ne(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ke(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(Ae(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function Pe(e,t,n){let r=new Map;return e.forEach(e=>{ke(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var Fe=F([J(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),J(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[V({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Ie=u({name:`Select`,props:Object.assign(Object.assign({},f.props),{to:m.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:o}=ee(e),s=f(`Select`,`-select`,Fe,ge,e,t),l=M(e.defaultValue),u=S(Y(e,`value`),l),p=M(!1),h=M(``),g=a(e,[`items`,`options`]),_=M([]),y=M([]),b=v(()=>y.value.concat(_.value).concat(g.value)),C=v(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return je(e,i);let a=t[r];return typeof a==`string`?je(e,a):typeof a==`number`?je(e,String(a)):!1}}),te=v(()=>{if(e.remote)return g.value;{let{value:t}=b,{value:n}=h;return!n.length||!e.filterable?t:Ne(t,C.value,n,e.childrenField)}}),ne=v(()=>{let{valueField:t,childrenField:n}=e,r=Me(t,n);return ue(te.value,r)}),w=v(()=>Pe(b.value,e.valueField,e.childrenField)),E=M(!1),D=S(Y(e,`show`),E),O=M(null),k=M(null),A=M(null),{localeRef:j}=d(`Select`),re=v(()=>e.placeholder??j.value.placeholder),N=[],P=M(new Map),ae=v(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function F(t){let n=e.remote,{value:r}=P,{value:i}=w,{value:a}=ae,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let L=v(()=>{if(e.multiple){let{value:e}=u;return Array.isArray(e)?F(e):[]}return null}),R=v(()=>{let{value:t}=u;return!e.multiple&&!Array.isArray(t)?t===null?null:F([t])[0]||null:null}),z=G(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:o?.value?.Select?.size||`medium`}}),{mergedSizeRef:V,mergedDisabledRef:H,mergedStatusRef:U}=z;function W(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=z;r&&I(r,t,n),a&&I(a,t,n),i&&I(i,t,n),l.value=t,o(),s()}function K(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=z;n&&I(n,t),r()}function oe(){let{onClear:t}=e;t&&I(t)}function q(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=z;n&&I(n,t),i(),r&&X()}function ce(t){let{onSearch:n}=e;n&&I(n,t)}function le(t){let{onScroll:n}=e;n&&I(n,t)}function J(){var t;let{remote:n,multiple:r}=e;if(n){let{value:n}=P;if(r){let{valueField:r}=e;(t=L.value)==null||t.forEach(e=>{n.set(e[r],e)})}else{let t=R.value;t&&n.set(t[e.valueField],t)}}}function de(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&I(n,t),r&&I(r,t),E.value=t}function X(){H.value||(de(!0),E.value=!0,e.filterable&&Le())}function Z(){de(!1)}function fe(){h.value=``,y.value=N}let pe=M(!1);function me(){e.filterable&&(pe.value=!0)}function he(){e.filterable&&(pe.value=!1,D.value||fe())}function _e(){H.value||(D.value?e.filterable?Le():Z():X())}function Q(e){(A.value?.selfRef)?.contains(e.relatedTarget)||(p.value=!1,K(e),Z())}function ve(e){q(e),p.value=!0}function ye(){p.value=!0}function be(e){O.value?.$el.contains(e.relatedTarget)||(p.value=!1,K(e),Z())}function $(){var e;(e=O.value)==null||e.focus(),Z()}function xe(e){D.value&&(O.value?.$el.contains(B(e))||Z())}function Se(t){if(!Array.isArray(t))return[];if(ae.value)return Array.from(t);{let{remote:n}=e,{value:r}=w;if(n){let{value:e}=P;return t.filter(t=>r.has(t)||e.has(t))}else return t.filter(e=>r.has(e))}}function Ce(e){we(e.rawNode)}function we(t){if(H.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=y,t=e[0]||null;if(t){let e=_.value;e.length?e.push(t):_.value=[t],y.value=N}}if(r&&P.value.set(t[a],t),e.multiple){let e=Se(u.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Te(t[a]);~e&&(_.value.splice(e,1),i&&(h.value=``))}}else e.push(t[a]),i&&(h.value=``);W(e,F(e))}else{if(n&&!r){let e=Te(t[a]);~e?_.value=[_.value[e]]:_.value=N}Ie(),Z(),W(t[a],t)}}function Te(t){return _.value.findIndex(n=>n[e.valueField]===t)}function Ee(t){D.value||X();let{value:n}=t.target;h.value=n;let{tag:r,remote:i}=e;if(ce(n),r&&!i){if(!n){y.value=N;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;g.value.some(e=>e[i]===r[i]||e[a]===r[a])||_.value.some(e=>e[i]===r[i]||e[a]===r[a])?y.value=N:y.value=[r]}}function De(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&Z(),r&&!i&&a&&(_.value=N),oe(),n?W([],[]):W(null,null)}function Oe(e){!se(e,`action`)&&!se(e,`empty`)&&!se(e,`header`)&&e.preventDefault()}function ke(e){le(e)}function Ae(t){var n,r,i;if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!O.value?.isComposing){if(D.value){let t=A.value?.getPendingTmNode();t?Ce(t):e.filterable||(Z(),Ie())}else if(X(),e.tag&&pe.value){let t=y.value[0];if(t){let n=t[e.valueField],{value:r}=u;e.multiple&&Array.isArray(r)&&r.includes(n)||we(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;D.value&&((n=A.value)==null||n.prev());break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;D.value?(r=A.value)==null||r.next():X();break;case`Escape`:D.value&&(ie(t),Z()),(i=O.value)==null||i.focus();break}}function Ie(){var e;(e=O.value)==null||e.focus()}function Le(){var e;(e=O.value)==null||e.focusInput()}function Re(){var e;D.value&&((e=k.value)==null||e.syncPosition())}J(),T(Y(e,`options`),J);let ze={focus:()=>{var e;(e=O.value)==null||e.focus()},focusInput:()=>{var e;(e=O.value)==null||e.focusInput()},blur:()=>{var e;(e=O.value)==null||e.blur()},blurInput:()=>{var e;(e=O.value)==null||e.blurInput()}},Be=v(()=>{let{self:{menuBoxShadow:e}}=s.value;return{"--n-menu-box-shadow":e}}),Ve=i?c(`select`,void 0,Be,e):void 0;return Object.assign(Object.assign({},ze),{mergedStatus:U,mergedClsPrefix:t,mergedBordered:n,namespace:r,treeMate:ne,isMounted:x(),triggerRef:O,menuRef:A,pattern:h,uncontrolledShow:E,mergedShow:D,adjustedTo:m(e),uncontrolledValue:l,mergedValue:u,followerRef:k,localizedPlaceholder:re,selectedOption:R,selectedOptions:L,mergedSize:V,mergedDisabled:H,focused:p,activeWithoutMenuOpen:pe,inlineThemeDisabled:i,onTriggerInputFocus:me,onTriggerInputBlur:he,handleTriggerOrMenuResize:Re,handleMenuFocus:ye,handleMenuBlur:be,handleMenuTabOut:$,handleTriggerClick:_e,handleToggle:Ce,handleDeleteOption:we,handlePatternInput:Ee,handleClear:De,handleTriggerBlur:Q,handleTriggerFocus:ve,handleKeydown:Ae,handleMenuAfterLeave:fe,handleMenuClickOutside:xe,handleMenuScroll:ke,handleMenuKeydown:Ae,handleMenuMousedown:Oe,mergedTheme:s,cssVars:i?void 0:Be,themeClass:Ve?.themeClass,onRender:Ve?.onRender})},render(){return W(`div`,{class:`${this.mergedClsPrefix}-select`},W(P,null,{default:()=>[W(h,null,{default:()=>W(Oe,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),W(X,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===m.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>W(e,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),E(W(Ee,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[p,this.mergedShow],[r,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[r,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{$ as i,Me as n,Ee as r,Ie as t};